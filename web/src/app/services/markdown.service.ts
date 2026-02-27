import { Injectable, SecurityContext } from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';
import hljs from 'highlight.js';
import { marked, Renderer } from 'marked';

@Injectable({
  providedIn: 'root'
})
export class MarkdownService {
  constructor(private sanitizer: DomSanitizer) {
    // Create custom renderer for syntax highlighting
    const renderer = new Renderer();
    renderer.code = ({ text, lang }: { text: string; lang?: string }) => {
      let highlightedCode = text;
      if (lang && hljs.getLanguage(lang)) {
        try {
          highlightedCode = hljs.highlight(text, { language: lang }).value;
        } catch (e) {
          highlightedCode = text;
        }
      } else {
        try {
          highlightedCode = hljs.highlightAuto(text).value;
        } catch (e) {
          highlightedCode = text;
        }
      }
      return `<pre><code class="hljs language-${lang || 'plaintext'}">${highlightedCode}</code></pre>`;
    };

    marked.setOptions({
      renderer,
      gfm: true,
    });
  }

  renderMarkdown(content: string): string {
    // Pre-process: Fix LLM output formatting issues
    let processedContent = (content || '');
    const rawHtml = marked.parse(processedContent);
    // Sanitize to prevent XSS
    return this.sanitizer.sanitize(SecurityContext.HTML, rawHtml) || '';
  }
}
