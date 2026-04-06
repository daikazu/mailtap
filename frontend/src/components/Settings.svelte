<script>
  import { onMount } from 'svelte';
  import { showSettings, smtpPort, themePref } from '../lib/stores';
  import { api } from '../lib/api';

  let port = '';
  let notifications = true;
  let theme = 'system';
  let needsRestart = false;
  let saved = false;

  onMount(async () => {
    const s = await api.getSettings();
    port = s.port;
    notifications = s.notifications;
    theme = s.theme;
  });

  async function save() {
    const current = await api.getSettings();
    needsRestart = port !== current.port;

    await api.saveSettings({ port, notifications, theme });

    // Apply theme immediately
    themePref.set(theme);
    localStorage.setItem('mailtap-theme', theme);

    if (!needsRestart) {
      smtpPort.set(port);
    }

    saved = true;
    setTimeout(() => { saved = false; }, 2000);
  }

  function close() {
    showSettings.set(false);
  }

  function handleKeydown(e) {
    if (e.key === 'Escape') {
      close();
    }
  }
</script>

<div class="command-palette-overlay" on:click={close} on:keydown={handleKeydown}>
  <div class="command-palette settings-modal" on:click|stopPropagation>
    <div class="settings-header">
      <h4>Settings</h4>
      <button class="hotkey-close" on:click={close}>&times;</button>
    </div>

    <div class="settings-body">
      <label class="setting-row">
        <span class="setting-label">SMTP Port</span>
        <input
          type="text"
          class="setting-input"
          bind:value={port}
          placeholder="2525"
        />
      </label>

      <label class="setting-row">
        <span class="setting-label">Notifications</span>
        <button
          class="toggle"
          class:active={notifications}
          on:click={() => notifications = !notifications}
        >
          <span class="toggle-knob"></span>
        </button>
      </label>

      <label class="setting-row">
        <span class="setting-label">Theme</span>
        <div class="theme-options">
          {#each ['system', 'dark', 'light'] as t}
            <button
              class="theme-btn"
              class:active={theme === t}
              on:click={() => theme = t}
            >{t}</button>
          {/each}
        </div>
      </label>
    </div>

    <div class="settings-footer">
      {#if needsRestart}
        <span class="restart-note">Restart app to change port</span>
      {/if}
      {#if saved && !needsRestart}
        <span class="saved-note">Saved</span>
      {/if}
      <button class="btn" on:click={save}>Save</button>
    </div>
  </div>
</div>

<style>
  .settings-modal {
    width: 400px;
  }

  .settings-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--border);
  }

  .settings-header h4 {
    margin: 0;
    font-size: 14px;
    color: var(--text-bright);
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

  .settings-body {
    padding: 16px 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .setting-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .setting-label {
    font-size: 13px;
    color: var(--text);
  }

  .setting-input {
    width: 80px;
    padding: 4px 8px;
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 4px;
    color: var(--text-bright);
    font-family: var(--font-mono);
    font-size: 12px;
    text-align: center;
    outline: none;
  }

  .setting-input:focus {
    border-color: var(--accent);
  }

  .toggle {
    width: 36px;
    height: 20px;
    border-radius: 10px;
    background: var(--border);
    border: none;
    cursor: pointer;
    position: relative;
    transition: background 0.2s;
    padding: 0;
  }

  .toggle.active {
    background: var(--accent);
  }

  .toggle-knob {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--bg-elevated);
    transition: transform 0.2s;
  }

  .toggle.active .toggle-knob {
    transform: translateX(16px);
  }

  .theme-options {
    display: flex;
    gap: 4px;
  }

  .theme-btn {
    padding: 4px 10px;
    font-size: 11px;
    font-family: var(--font-mono);
    background: var(--bg);
    border: 1px solid var(--border);
    border-radius: 4px;
    color: var(--text-dim);
    cursor: pointer;
    text-transform: capitalize;
    transition: border-color 0.15s, color 0.15s;
  }

  .theme-btn.active {
    border-color: var(--accent);
    color: var(--accent);
  }

  .settings-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 12px;
    padding: 12px 20px;
    border-top: 1px solid var(--border);
  }

  .restart-note {
    font-size: 11px;
    color: var(--accent);
  }

  .saved-note {
    font-size: 11px;
    color: var(--accent);
  }
</style>
