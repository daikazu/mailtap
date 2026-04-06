<script>
  import { onMount } from 'svelte';
  import { showHotkeySheet } from '../lib/stores';
  import { api } from '../lib/api';

  let port = '2525';

  onMount(async () => {
    port = await api.getPort();
  });

  $: envSnippet = `MAIL_MAILER=smtp
MAIL_HOST=127.0.0.1
MAIL_PORT=${port}
MAIL_USERNAME=email@example.com
MAIL_ENCRYPTION=null`;

  let copied = false;

  function copyEnv() {
    navigator.clipboard.writeText(envSnippet);
    copied = true;
    setTimeout(() => copied = false, 2000);
  }

  const hotkeys = [
    { action: 'Command palette', keys: ['Ctrl+K'] },
    { action: 'Search', keys: ['/'] },
    { action: 'Next email', keys: ['j', '\u2193'] },
    { action: 'Previous email', keys: ['k', '\u2191'] },
    { action: 'Open email', keys: ['Enter'] },
    { action: 'Close or Back', keys: ['Esc'] },
    { action: 'Delete email', keys: ['d', 'Del'] },
    { action: 'Previous tab', keys: ['h', '\u2190'] },
    { action: 'Next tab', keys: ['l', '\u2192'] },
    { action: 'Jump to tab', keys: ['1-5'] },
    { action: 'Clear all', keys: ['Ctrl+Shift+Del'] },
    { action: 'Settings', keys: [','] },
    { action: 'This help', keys: ['?'] },
  ];

  function handleClickOutside(e) {
    if (e.target === e.currentTarget) {
      showHotkeySheet.set(false);
    }
  }
</script>

<div class="hotkey-sheet-overlay" on:click={handleClickOutside}>
  <div class="hotkey-sheet">
    <div class="hotkey-sheet-header">
      <h4>Keyboard Shortcuts</h4>
      <button class="hotkey-close" on:click={() => showHotkeySheet.set(false)}>&times;</button>
    </div>
    {#each hotkeys as hk}
      <div class="hotkey-row">
        <span class="action">{hk.action}</span>
        <span class="keys">
          {#each hk.keys as key, i}
            {#if i > 0}<span class="or">or</span>{/if}
            <kbd>{key}</kbd>
          {/each}
        </span>
      </div>
    {/each}

    <div class="config-section">
      <div class="config-header">
        <h4>.env Configuration</h4>
        <button class="copy-btn" on:click={copyEnv}>
          {copied ? 'Copied!' : 'Copy'}
        </button>
      </div>
      <pre class="env-block">{envSnippet}</pre>
    </div>
  </div>
</div>

<style>
  .hotkey-sheet-overlay {
    position: fixed;
    inset: 0;
    z-index: 100;
  }

  .hotkey-sheet {
    position: fixed;
    bottom: 40px;
    right: 16px;
    top: auto;
    left: auto;
    transform: none;
  }

  .hotkey-sheet-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }

  .hotkey-sheet-header h4 {
    margin-bottom: 0;
  }

  .keys {
    display: inline-flex;
    align-items: center;
    gap: 4px;
  }

  .keys kbd {
    color: var(--accent);
    background: var(--accent-dim);
    border-color: var(--accent);
    margin-left: 0;
  }

  .or {
    font-size: 10px;
    color: var(--text-dim);
    position: relative;
    top: 1px;
  }

  .hotkey-close {
    background: none;
    border: none;
    color: var(--text-dim);
    font-size: 18px;
    cursor: pointer;
    padding: 0 4px;
    font-family: var(--font-mono);
    line-height: 1;
  }

  .hotkey-close:hover {
    color: var(--text);
  }

  .config-section {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid var(--border);
  }

  .config-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .config-header h4 {
    margin-bottom: 0;
  }

  .copy-btn {
    background: var(--bg);
    border: 1px solid var(--border);
    color: var(--text-dim);
    font-family: var(--font-mono);
    font-size: 10px;
    padding: 2px 8px;
    border-radius: 3px;
    cursor: pointer;
    transition: border-color 0.15s, color 0.15s;
  }

  .copy-btn:hover {
    border-color: var(--accent);
    color: var(--accent);
  }

  .env-block {
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--accent);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 10px 12px;
    margin: 0;
    line-height: 1.6;
    white-space: pre;
    user-select: text;
    -webkit-user-select: text;
  }
</style>
