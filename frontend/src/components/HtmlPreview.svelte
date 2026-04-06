<script>
  import { onMount, onDestroy } from 'svelte';

  export let html = '';

  // Forward keyboard events from sandboxed iframe via postMessage
  const forwardScript = `<script>
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

  function handleMessage(e) {
    if (e.data?.type === 'mailtap-keydown') {
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
      srcdoc={srcdoc}
      sandbox="allow-scripts"
      style="width: 100%; flex: 1; border: none; background: #ffffff; min-height: 0;"
      title="HTML Preview"
    ></iframe>
  {:else}
    <p style="color: var(--text-dim);">No HTML content</p>
  {/if}
</div>
