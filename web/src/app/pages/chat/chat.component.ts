import { CommonModule } from '@angular/common';
import { AfterViewChecked, ChangeDetectionStrategy, Component, ElementRef, ViewChild, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { Message } from '../../models/config.model';
import { ApiService } from '../../services/api.service';
import { MarkdownService } from '../../services/markdown.service';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatProgressBarModule,
    MatTooltipModule
  ],
  templateUrl: './chat.component.html',
  host: { style: 'display: block; height: 100%; min-height: 0;' },
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ChatComponent implements AfterViewChecked {
  @ViewChild('messagesContainer') private messagesContainer!: ElementRef;

  messages = signal<Message[]>([]);
  inputMessage = signal('');
  isLoading = signal(false);
  private lastScrollKey = '';

  readonly quickPrompts = [
    'What can you help me with?',
    'Explain quantum computing in simple terms',
    'Write a Python function to calculate Fibonacci',
    'Help me debug my code'
  ];

  constructor(
    private apiService: ApiService,
    private markdownService: MarkdownService
  ) {}

  renderMarkdown(content: string): string { return this.markdownService.renderMarkdown(content); }

  toggleRawContent(message: Message): void {
    message.showRawContent = !message.showRawContent;
  }

  ngAfterViewChecked(): void {
  }

  onKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      this.sendMessage();
    }
  }

  sendQuickMessage(content: string): void {
    this.inputMessage.set(content);
    this.sendMessage();
  }

  sendMessage(): void {
    const content = this.inputMessage().trim();
    if (!content || this.isLoading()) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content,
      timestamp: new Date(),
        isComplete: true
    };

    this.messages.update(msgs => [...msgs, userMessage]);
    this.inputMessage.set('');
    this.isLoading.set(true);

    // For non-streaming API: initialize assistant message as complete (not streaming)
    const assistantMessage: Message = {
      id: (Date.now() + 1).toString(),
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      isComplete: false
    };

    this.messages.update(msgs => [...msgs, assistantMessage]);

    // sendMessage returns a single response object: { response: string }
    this.apiService.sendMessage(content).subscribe({
      next: (response: any) => {
        const text = response?.response ?? 'No response';
        this.messages.update(msgs => {
          const lastMsg = msgs[msgs.length - 1];
          if (lastMsg.role === 'assistant') {
            return [
              ...msgs.slice(0, -1),
              { ...lastMsg, content: text, isStreaming: false, isComplete: true }
            ];
          }
          return msgs;
        });
        this.isLoading.set(false);
      },
      error: (error: any) => {
        this.messages.update(msgs => {
          const lastMsg = msgs[msgs.length - 1];
          if (lastMsg.role === 'assistant') {
            return [
              ...msgs.slice(0, -1),
              { ...lastMsg, content: `Error: ${error?.message ?? String(error)}`, isStreaming: false, isComplete: true }
            ];
          }
          return msgs;
        });
        this.isLoading.set(false);
      }
    });
  }
}
