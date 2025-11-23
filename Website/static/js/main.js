// Main JavaScript for Universal Ebook Translator Website

// Navigation functionality
class Navigation {
  constructor() {
    this.navToggle = document.querySelector('.mobile-menu-toggle');
    this.navLinks = document.querySelector('.nav-links');
    this.init();
  }

  init() {
    if (this.navToggle) {
      this.navToggle.addEventListener('click', () => this.toggleMobileMenu());
    }

    // Close mobile menu when clicking outside
    document.addEventListener('click', (e) => {
      if (this.navLinks && this.navLinks.classList.contains('active') && 
          !this.navToggle.contains(e.target) && 
          !this.navLinks.contains(e.target)) {
        this.closeMobileMenu();
      }
    });

    // Active navigation highlighting
    this.updateActiveNav();
    window.addEventListener('scroll', () => this.updateActiveNav());
  }

  toggleMobileMenu() {
    if (this.navLinks) {
      this.navLinks.classList.toggle('active');
      this.navToggle.classList.toggle('active');
    }
  }

  closeMobileMenu() {
    if (this.navLinks) {
      this.navLinks.classList.remove('active');
      this.navToggle.classList.remove('active');
    }
  }

  updateActiveNav() {
    const sections = document.querySelectorAll('main > section');
    const navLinks = document.querySelectorAll('.nav-link');
    
    if (sections.length === 0 || navLinks.length === 0) return;

    let currentSection = '';
    sections.forEach(section => {
      const sectionTop = section.offsetTop;
      const sectionHeight = section.offsetHeight;
      if (window.pageYOffset >= sectionTop - 100) {
        currentSection = section.getAttribute('id');
      }
    });

    navLinks.forEach(link => {
      link.classList.remove('active');
      if (link.getAttribute('href') === `#${currentSection}`) {
        link.classList.add('active');
      }
    });
  }
}

// API Playground functionality
class APIPlayground {
  constructor() {
    this.container = document.querySelector('.demo-container');
    this.form = null;
    this.resultArea = null;
    this.init();
  }

  init() {
    if (!this.container) return;

    this.createPlayground();
    this.bindEvents();
  }

  createPlayground() {
    this.container.innerHTML = `
      <h3>Try the Translation API</h3>
      <p>Translate text using our advanced translation engine. Limited to 500 characters for demo.</p>
      
      <form class="demo-form" id="translationForm">
        <div class="form-group">
          <label class="form-label" for="sourceText">Source Text</label>
          <textarea 
            class="form-textarea" 
            id="sourceText" 
            name="sourceText"
            placeholder="Enter text to translate..."
            maxlength="500"
            required
          ></textarea>
          <small id="charCount">0 / 500 characters</small>
        </div>
        
        <div class="form-group">
          <label class="form-label" for="sourceLang">Source Language</label>
          <select class="form-select" id="sourceLang" name="sourceLang" required>
            <option value="">Select source language</option>
            <option value="en">English</option>
            <option value="fr">French</option>
            <option value="de">German</option>
            <option value="es">Spanish</option>
            <option value="it">Italian</option>
            <option value="pt">Portuguese</option>
            <option value="ru">Russian</option>
            <option value="zh">Chinese</option>
            <option value="ja">Japanese</option>
            <option value="ko">Korean</option>
          </select>
        </div>
        
        <div class="form-group">
          <label class="form-label" for="targetLang">Target Language</label>
          <select class="form-select" id="targetLang" name="targetLang" required>
            <option value="">Select target language</option>
            <option value="en">English</option>
            <option value="fr">French</option>
            <option value="de">German</option>
            <option value="es">Spanish</option>
            <option value="it">Italian</option>
            <option value="pt">Portuguese</option>
            <option value="ru">Russian</option>
            <option value="zh">Chinese</option>
            <option value="ja">Japanese</option>
            <option value="ko">Korean</option>
          </select>
        </div>
        
        <div class="form-group">
          <label class="form-label" for="provider">Translation Provider</label>
          <select class="form-select" id="provider" name="provider" required>
            <option value="openai">OpenAI GPT-4</option>
            <option value="anthropic">Anthropic Claude</option>
            <option value="zhipu">Zhipu AI GLM-4</option>
            <option value="deepseek">DeepSeek</option>
            <option value="qwen">Qwen</option>
            <option value="gemini">Google Gemini</option>
          </select>
        </div>
        
        <button type="submit" class="btn btn-primary">
          <span class="btn-text">Translate</span>
          <span class="btn-spinner" style="display: none;">Translating...</span>
        </button>
      </form>
      
      <div id="translationResult" style="display: none; margin-top: 2rem;">
        <h4>Translation Result</h4>
        <div id="resultContent" style="background: var(--background-tertiary); padding: 1rem; border-radius: var(--border-radius); margin-top: 1rem;">
          <!-- Result will be inserted here -->
        </div>
      </div>
    `;

    this.form = document.getElementById('translationForm');
    this.resultArea = document.getElementById('translationResult');
  }

  bindEvents() {
    if (!this.form) return;

    // Character counter
    const sourceText = document.getElementById('sourceText');
    const charCount = document.getElementById('charCount');
    
    if (sourceText && charCount) {
      sourceText.addEventListener('input', () => {
        const length = sourceText.value.length;
        charCount.textContent = `${length} / 500 characters`;
        charCount.style.color = length > 450 ? 'var(--warning-color)' : 'var(--text-secondary)';
      });
    }

    // Form submission
    this.form.addEventListener('submit', (e) => {
      e.preventDefault();
      this.translate();
    });

    // Language auto-swap
    const sourceLang = document.getElementById('sourceLang');
    const targetLang = document.getElementById('targetLang');
    
    if (sourceLang && targetLang) {
      const autoSwap = () => {
        if (sourceLang.value && targetLang.value && sourceLang.value === targetLang.value) {
          // Find first different language
          const options = Array.from(targetLang.options);
          const differentOption = options.find(option => 
            option.value && option.value !== sourceLang.value
          );
          if (differentOption) {
            targetLang.value = differentOption.value;
          }
        }
      };

      sourceLang.addEventListener('change', autoSwap);
      targetLang.addEventListener('change', autoSwap);
    }
  }

  async translate() {
    if (!this.form) return;

    const formData = new FormData(this.form);
    const data = Object.fromEntries(formData);

    const btnText = this.form.querySelector('.btn-text');
    const btnSpinner = this.form.querySelector('.btn-spinner');

    // Show loading state
    if (btnText) btnText.style.display = 'none';
    if (btnSpinner) btnSpinner.style.display = 'inline';

    try {
      // Simulate API call (replace with actual API endpoint)
      const response = await this.mockTranslationAPI(data);
      
      // Show result
      this.showResult(response);
    } catch (error) {
      this.showError(error.message);
    } finally {
      // Hide loading state
      if (btnText) btnText.style.display = 'inline';
      if (btnSpinner) btnSpinner.style.display = 'none';
    }
  }

  async mockTranslationAPI(data) {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 1500));

    // Mock translation (in real implementation, this would call the actual API)
    const mockTranslations = {
      'en-fr': {
        'hello': 'bonjour',
        'goodbye': 'au revoir',
        'thank you': 'merci',
        'how are you': 'comment allez-vous'
      },
      'en-de': {
        'hello': 'hallo',
        'goodbye': 'auf wiedersehen',
        'thank you': 'danke',
        'how are you': 'wie geht es dir'
      },
      'en-es': {
        'hello': 'hola',
        'goodbye': 'adiós',
        'thank you': 'gracias',
        'how are you': 'cómo estás'
      }
    };

    const langPair = `${data.sourceLang}-${data.targetLang}`;
    const translations = mockTranslations[langPair] || {};
    const lowerText = data.sourceText.toLowerCase().trim();
    
    let translatedText = data.sourceText; // Default to original if no translation found
    
    // Simple mock translation for common phrases
    for (const [english, translated] of Object.entries(translations)) {
      if (lowerText.includes(english)) {
        translatedText = data.sourceText.replace(english, translated);
        break;
      }
    }

    return {
      original: data.sourceText,
      translated: translatedText,
      sourceLang: data.sourceLang,
      targetLang: data.targetLang,
      provider: data.provider,
      confidence: 0.85 + Math.random() * 0.15, // Mock confidence score
      processingTime: (1 + Math.random() * 2).toFixed(1) + 's'
    };
  }

  showResult(result) {
    if (!this.resultArea) return;

    const resultContent = document.getElementById('resultContent');
    if (!resultContent) return;

    resultContent.innerHTML = `
      <div style="margin-bottom: 1rem;">
        <strong>Original (${result.sourceLang}):</strong>
        <p style="margin: 0.5rem 0; padding: 0.75rem; background: var(--background); border-radius: var(--border-radius);">
          ${result.original}
        </p>
      </div>
      <div style="margin-bottom: 1rem;">
        <strong>Translated (${result.targetLang}):</strong>
        <p style="margin: 0.5rem 0; padding: 0.75rem; background: var(--background); border-radius: var(--border-radius); border-left: 4px solid var(--success-color);">
          ${result.translated}
        </p>
      </div>
      <div style="display: flex; gap: 1rem; font-size: 0.875rem; color: var(--text-secondary);">
        <span>Provider: <strong>${result.provider}</strong></span>
        <span>Confidence: <strong>${(result.confidence * 100).toFixed(1)}%</strong></span>
        <span>Processing time: <strong>${result.processingTime}</strong></span>
      </div>
    `;

    this.resultArea.style.display = 'block';
    this.resultArea.scrollIntoView({ behavior: 'smooth' });
  }

  showError(message) {
    if (!this.resultArea) return;

    const resultContent = document.getElementById('resultContent');
    if (!resultContent) return;

    resultContent.innerHTML = `
      <div style="padding: 1rem; background: var(--error-color); color: white; border-radius: var(--border-radius);">
        <strong>Error:</strong> ${message}
      </div>
    `;

    this.resultArea.style.display = 'block';
    this.resultArea.scrollIntoView({ behavior: 'smooth' });
  }
}

// Documentation enhancements
class Documentation {
  constructor() {
    this.init();
  }

  init() {
    this.addCopyButtons();
    this.addCodeLineNumbers();
    this.addTableOfContents();
    this.addSearch();
  }

  addCopyButtons() {
    const codeBlocks = document.querySelectorAll('pre code');
    
    codeBlocks.forEach(block => {
      const button = document.createElement('button');
      button.className = 'copy-button';
      button.textContent = 'Copy';
      button.style.cssText = `
        position: absolute;
        top: 0.5rem;
        right: 0.5rem;
        padding: 0.25rem 0.5rem;
        font-size: 0.75rem;
        background: var(--primary-color);
        color: white;
        border: none;
        border-radius: var(--border-radius);
        cursor: pointer;
        opacity: 0;
        transition: opacity 0.2s;
      `;

      const pre = block.parentElement;
      pre.style.position = 'relative';
      pre.appendChild(button);

      pre.addEventListener('mouseenter', () => button.style.opacity = '1');
      pre.addEventListener('mouseleave', () => button.style.opacity = '0');

      button.addEventListener('click', () => {
        navigator.clipboard.writeText(block.textContent).then(() => {
          button.textContent = 'Copied!';
          setTimeout(() => button.textContent = 'Copy', 2000);
        });
      });
    });
  }

  addCodeLineNumbers() {
    const codeBlocks = document.querySelectorAll('pre code');
    
    codeBlocks.forEach(block => {
      const lines = block.textContent.split('\n');
      const lineNumbers = lines.map((_, i) => i + 1).join('\n');
      
      const lineNumbersElement = document.createElement('span');
      lineNumbersElement.textContent = lineNumbers;
      lineNumbersElement.style.cssText = `
        display: inline-block;
        width: 2rem;
        text-align: right;
        margin-right: 1rem;
        color: var(--text-light);
        user-select: none;
        border-right: 1px solid var(--border-color);
        padding-right: 0.5rem;
        margin-right: 0.5rem;
      `;

      block.style.display = 'flex';
      block.style.alignItems = 'flex-start';
      block.insertBefore(lineNumbersElement, block.firstChild);
    });
  }

  addTableOfContents() {
    const content = document.querySelector('.docs-content');
    if (!content) return;

    const headings = content.querySelectorAll('h2, h3, h4');
    if (headings.length === 0) return;

    const toc = document.createElement('div');
    toc.className = 'table-of-contents';
    toc.innerHTML = '<h4>Table of Contents</h4><ul></ul>';
    
    const tocList = toc.querySelector('ul');
    
    headings.forEach(heading => {
      const id = heading.textContent.toLowerCase().replace(/\s+/g, '-').replace(/[^\w-]/g, '');
      heading.id = id;
      
      const li = document.createElement('li');
      li.style.marginLeft = `${(parseInt(heading.tagName.charAt(1)) - 2) * 1rem}rem`;
      
      const link = document.createElement('a');
      link.href = `#${id}`;
      link.textContent = heading.textContent;
      link.style.cssText = `
        color: var(--text-secondary);
        text-decoration: none;
        font-size: 0.875rem;
        line-height: 1.5;
        transition: color 0.2s;
      `;
      
      link.addEventListener('mouseenter', () => link.style.color = 'var(--primary-color)');
      link.addEventListener('mouseleave', () => link.style.color = 'var(--text-secondary)');
      
      li.appendChild(link);
      tocList.appendChild(li);
    });

    content.insertBefore(toc, content.firstChild);
  }

  addSearch() {
    const searchContainer = document.querySelector('.search-container');
    if (!searchContainer) return;

    const input = document.createElement('input');
    input.type = 'text';
    input.placeholder = 'Search documentation...';
    input.className = 'search-input';
    input.style.cssText = `
      width: 100%;
      padding: 0.75rem;
      border: 1px solid var(--border-color);
      border-radius: var(--border-radius);
      font-size: 1rem;
      margin-bottom: 1rem;
    `;

    searchContainer.appendChild(input);

    input.addEventListener('input', (e) => {
      const query = e.target.value.toLowerCase();
      const sections = document.querySelectorAll('.docs-content section');
      
      sections.forEach(section => {
        const text = section.textContent.toLowerCase();
        if (text.includes(query) || query === '') {
          section.style.display = 'block';
        } else {
          section.style.display = 'none';
        }
      });
    });
  }
}

// Theme toggle
class ThemeManager {
  constructor() {
    this.themeToggle = document.querySelector('.theme-toggle');
    this.init();
  }

  init() {
    if (!this.themeToggle) return;

    // Load saved theme
    const savedTheme = localStorage.getItem('theme') || 'light';
    this.setTheme(savedTheme);

    this.themeToggle.addEventListener('click', () => {
      const currentTheme = document.body.getAttribute('data-theme') || 'light';
      const newTheme = currentTheme === 'light' ? 'dark' : 'light';
      this.setTheme(newTheme);
    });
  }

  setTheme(theme) {
    document.body.setAttribute('data-theme', theme);
    localStorage.setItem('theme', theme);
    
    if (this.themeToggle) {
      this.themeToggle.textContent = theme === 'light' ? '🌙' : '☀️';
    }
  }
}

// Analytics and performance monitoring
class Analytics {
  constructor() {
    this.init();
  }

  init() {
    // Simple page view tracking
    this.trackPageView();
    
    // Track external link clicks
    this.trackExternalLinks();
    
    // Track scroll depth
    this.trackScrollDepth();
  }

  trackPageView() {
    const page = window.location.pathname;
    console.log('Page view:', page);
    // In real implementation, send to analytics service
  }

  trackExternalLinks() {
    const externalLinks = document.querySelectorAll('a[href^="http"]');
    
    externalLinks.forEach(link => {
      link.addEventListener('click', () => {
        const url = link.getAttribute('href');
        console.log('External link clicked:', url);
        // In real implementation, send to analytics service
      });
    });
  }

  trackScrollDepth() {
    let maxScroll = 0;
    
    window.addEventListener('scroll', () => {
      const scrollPercentage = (window.scrollY / (document.body.scrollHeight - window.innerHeight)) * 100;
      maxScroll = Math.max(maxScroll, scrollPercentage);
    });

    // Track scroll depth on page unload
    window.addEventListener('beforeunload', () => {
      if (maxScroll > 50) {
        console.log('Scroll depth:', Math.round(maxScroll) + '%');
        // In real implementation, send to analytics service
      }
    });
  }
}

// Initialize everything when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new Navigation();
  new APIPlayground();
  new Documentation();
  new ThemeManager();
  new Analytics();
});

// Performance optimizations
window.addEventListener('load', () => {
  // Preload critical resources
  const criticalResources = [
    '/static/css/main.css',
    '/static/js/main.js'
  ];
  
  criticalResources.forEach(resource => {
    const link = document.createElement('link');
    link.rel = 'preload';
    link.as = resource.endsWith('.css') ? 'style' : 'script';
    link.href = resource;
    document.head.appendChild(link);
  });
});

// Service Worker registration for offline support
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js')
      .then(registration => {
        console.log('SW registered:', registration);
      })
      .catch(error => {
        console.log('SW registration failed:', error);
      });
  });
}