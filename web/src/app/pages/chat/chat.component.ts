import { Component, signal, ViewChild, ElementRef, AfterViewChecked } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatListModule } from '@angular/material/list';
import { MatDividerModule } from '@angular/material/divider';
import { ApiService } from '../../services/api.service';
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
    MatTooltipModule,
    MatListModule,
    MatDividerModule
  ],
  template: `
    <div class="chat-wrapper">
      <!-- Messages Area -->
      <div #messagesContainer class="messages-area">
        @if (messages().length > 0) {
          <mat-list class="messages-list">
            @for (message of messages(); track message.id) {
              <mat-list-item class="message-item" [class.user-message]="message.role === 'user'">
                <mat-icon matListItemIcon>
                  {{ message.role === 'user' ? 'person' : 'smart_toy' }}
                </mat-icon>
                <div matListItemTitle class="message-header">
                  <span>{{ message.role === 'user' ? 'You' : 'PicoClaw' }}</span>
                  <span class="timestamp">{{ message.timestamp | date:'shortTime' }}</span>
                </div>
                <div matListItemLine class="message-content">
                  <pre>{{ message.content }}</pre>
                  @if (message.isStreaming) {
                    <span class="streaming-indicator">|</span>
                  }
                </div>
              </mat-list-item>
              <mat-divider></mat-divider>
            }
            
            @if (isLoading() && !streamingMessage()) {
              <mat-list-item class="thinking-item">
                <mat-spinner matListItemIcon diameter="20"></mat-spinner>
                <div matListItemTitle>Thinking...</div>
              </mat-list-item>
            }
          </mat-list>
        }
      </div>

      <!-- Input Area - Fixed at bottom -->
      <div class="input-area">
        <mat-divider></mat-divider>
        <div class="input-container">
          <mat-form-field appearance="outline" class="message-input">
            <mat-label>Type your message...</mat-label>
            <textarea matInput
                      [(ngModel)]="inputMessage"
                      (keydown)="onKeydown($event)"
                      cdkTextareaAutosize
                      cdkAutosizeMinRows="1"
                      cdkAutosizeMaxRows="5">
            </textarea>
          </mat-form-field>
          
          <button mat-fab
                  color="primary"
                  [disabled]="!inputMessage().trim() || isLoading()"
                  (click)="sendMessage()"
                  matTooltip="Send message (Enter)">
            @if (isLoading()) {
              <mat-spinner diameter="24"></mat-spinner>
            } @else {
              <mat-icon>send</mat-icon>
            }
          </button>
        </div>
      </div>
    </div>
  `,
  styles: [`
    :host {
      display: block;
      height: 100%;
    }

    .chat-wrapper {
      display: flex;
      flex-direction: column;
      height: 100%;
      max-width: 800px;
      margin: 0 auto;
    }

    .messages-area {
      flex: 1;
      overflow-y: auto;
      padding: 16px;
    }

    .messages-list {
      padding: 0;
    }

    .message-item {
      height: auto !important;
      min-height: 72px;
      padding: 16px 0;
    }

    .user-message {
      background-color: var(--mat-sys-primary-container);
    }

    .message-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;
    }

    .timestamp {
      font-size: 12px;
      opacity: 0.6;
    }

    .message-content pre {
      margin: 0;
      white-space: pre-wrap;
      word-wrap: break-word;
      font-family: inherit;
      font-size: 14px;
    }

    .streaming-indicator {
      animation: blink 1s infinite;
    }

    @keyframes blink {
      0%, 50% { opacity: 1; }
      51%, 100% { opacity: 0; }
    }

    .thinking-item {
      color: var(--mat-sys-on-surface-variant);
    }

    .input-area {
      flex-shrink: 0;
      background: var(--mat-sys-surface);
    }

    .input-container {
      display: flex;
      gap: 12px;
      align-items: flex-start;
      padding: 16px;
    }

    .message-input {
      flex: 1;
    }
  `]
})
export class ChatComponent implements AfterViewChecked {
  @ViewChild('messagesContainer') private messagesContainer!: ElementRef;

  messages = signal<Message[]>([]);
  inputMessage = signal('');
  isLoading = signal(false);
  streamingMessage = signal('');

  quickPrompts = [
    'What can you help me with?',
    'Explain quantum computing in simple terms',
    'Write a Python function to calculate Fibonacci',
    'Help me debug my code'
  ];

  constructor(private apiService: ApiService) {}

  ngAfterViewChecked(): void {
    this.scrollToBottom();
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
      isStreaming: true
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
              { ...lastMsg, isStreaming: false }
            ];
          }
          return msgs;
        });
        this.isLoading.set(false);
      }
    });
  }
}
