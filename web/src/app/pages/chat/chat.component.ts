import { AfterViewChecked, ChangeDetectionStrategy, Component, ElementRef, ViewChild, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ApiService } from '../../services/api.service';
import { MarkdownService } from '../../services/markdown.service';
import { Message } from '../../models/config.model';

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

  toggleReasoning(message: Message): void {
    message.isReasoningCollapsed = !message.isReasoningCollapsed;
  }

  ngAfterViewChecked(): void {
    const message = this.messages().at(-1);
    const scrollKey = message ? `${message.id}:${message.content.length}:${message.isStreaming ? '1' : '0'}` : 'empty';
    if (scrollKey !== this.lastScrollKey) {
      this.lastScrollKey = scrollKey;
      this.scrollToBottom();
    }
  }

  private scrollToBottom(): void {
    try {
      const element = this.messagesContainer?.nativeElement;
      if (element) {
        element.scrollTop = element.scrollHeight;
      }
    } catch (err) {
      console.error('Error scrolling to bottom:', err);
    }
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
      timestamp: new Date()
    };

    this.messages.update(msgs => [...msgs, userMessage]);
    this.inputMessage.set('');
    this.isLoading.set(true);

    let fullResponse = '';
    const assistantMessage: Message = {
      id: (Date.now() + 1).toString(),
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      isStreaming: true,
      isComplete: false
    };

    this.messages.update(msgs => [...msgs, assistantMessage]);

    this.apiService.streamMessage(content).subscribe({
      next: (chunk) => {
        fullResponse += chunk;
        this.messages.update(msgs => {
          const lastMsg = msgs[msgs.length - 1];
          if (lastMsg.role === 'assistant') {
            return [
              ...msgs.slice(0, -1),
              { ...lastMsg, content: fullResponse }
            ];
          }
          return msgs;
        });
      },
      error: () => {
        this.apiService.sendMessage(content).subscribe({
          next: (response) => {
            this.messages.update(msgs => {
              const lastMsg = msgs[msgs.length - 1];
              if (lastMsg.role === 'assistant') {
                return [
                  ...msgs.slice(0, -1),
                  { ...lastMsg, content: response.response || 'No response', isStreaming: false }
                ];
              }
              return msgs;
            });
            this.isLoading.set(false);
          },
          error: (error) => {
            this.messages.update(msgs => {
              const lastMsg = msgs[msgs.length - 1];
              if (lastMsg.role === 'assistant') {
                return [
                  ...msgs.slice(0, -1),
                  { ...lastMsg, content: `Error: ${error.message}`, isStreaming: false }
                ];
              }
              return msgs;
            });
            this.isLoading.set(false);
          }
        });
      },
      complete: () => {
        this.messages.update(msgs => {
          const lastMsg = msgs[msgs.length - 1];
          if (lastMsg.role === 'assistant') {
            return [
              ...msgs.slice(0, -1),
              { ...lastMsg, isStreaming: false, isComplete: true }
            ];
          }
          return msgs;
        });
        this.isLoading.set(false);
      }
    });
  }
}
