<script>
  import { onMount, onDestroy } from 'svelte';
  import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';

  export let html = '';

  let iframe;

  // Forward keyboard events and link clicks from sandboxed iframe via postMessage
  const forwardScript = `<script>
    document.addEventListener('click', function(e) {
      var a = e.target.closest && e.target.closest('a[href]');
      if (!a) return;
      var href = a.getAttribute('href');
      if (!href || href.charAt(0) === '#') return;
      e.preventDefault();
      window.parent.postMessage({ type: 'mailtap-open-link', href: a.href }, '*');
    });
    document.addEventListener('keydown', function(e) {
      window.parent.postMessage({
        type: 'mailtap-keydown',
        key: e.key, code: e.code, keyCode: e.keyCode,
        ctrlKey: e.ctrlKey, shiftKey: e.shiftKey, altKey: e.altKey, metaKey: e.metaKey
      }, '*');
    });
  <\/script>`;

  function injectForwardScript(src) {
    if (!src) return '';
    // Case-insensitive </body> replacement
    const bodyClose = src.match(/<\/body\s*>/i);
    if (bodyClose) {
      return src.replace(bodyClose[0], forwardScript + bodyClose[0]);
    }
    return src + forwardScript;
  }

  $: srcdoc = html ? injectForwardScript(html) : '';

  const allowedProtocols = ['http:', 'https:', 'mailto:'];

  function openLink(href) {
    let url;
    try {
      url = new URL(href);
    } catch {
      return;
    }
    if (allowedProtocols.includes(url.protocol)) {
      BrowserOpenURL(url.href);
    }
  }

  function handleMessage(e) {
    if (!iframe || e.source !== iframe.contentWindow) return;
    if (e.data?.type === 'mailtap-open-link') {
      openLink(e.data.href);
    } else if (e.data?.type === 'mailtap-keydown') {
      window.dispatchEvent(new KeyboardEvent('keydown', {
        key: e.data.key, code: e.data.code, keyCode: e.data.keyCode,
        ctrlKey: e.data.ctrlKey, shiftKey: e.data.shiftKey,
        altKey: e.data.altKey, metaKey: e.data.metaKey,
        bubbles: true
      }));
    }
  }

  onMount(() => window.addEventListener('message', handleMessage));
  onDestroy(() => window.removeEventListener('message', handleMessage));
</script>

<div class="html-preview" style="display: flex; flex-direction: column; flex: 1; height: 100%;">
  {#if html}
    <iframe
      bind:this={iframe}
      srcdoc={srcdoc}
      sandbox="allow-scripts"
      style="width: 100%; flex: 1; border: none; background: #ffffff; min-height: 0;"
      title="HTML Preview"
    ></iframe>
  {:else}
    <p style="color: var(--text-dim);">No HTML content</p>
  {/if}
</div>
