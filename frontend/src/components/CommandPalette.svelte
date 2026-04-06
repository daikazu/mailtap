<script>
  import { onMount, tick } from 'svelte';
  import { get } from 'svelte/store';
  import { showCommandPalette, activeTab, showHotkeySheet, showConfirmClear, showSettings, selectedEmailId, deleteEmail, themePref } from '../lib/stores';

  let query = '';
  let selectedIndex = 0;
  let inputEl;

  const commands = [
    { label: 'Delete Email', key: 'd', action: () => { const id = get(selectedEmailId); if (id) deleteEmail(id); } },
    { label: 'Delete All Emails', key: 'Ctrl+Shift+Del', action: () => { showConfirmClear.set(true); } },
    { label: 'HTML Preview', key: '1', action: () => activeTab.set(0) },
    { label: 'Plain Text', key: '2', action: () => activeTab.set(1) },
    { label: 'Headers', key: '3', action: () => activeTab.set(2) },
    { label: 'Attachments', key: '4', action: () => activeTab.set(3) },
    { label: 'Source', key: '5', action: () => activeTab.set(4) },
    { label: 'Focus Search', key: '/', action: () => document.querySelector('.search')?.focus() },
    { label: 'Settings', key: ',', action: () => showSettings.set(true) },
    { label: 'Show Hotkeys', key: '?', action: () => showHotkeySheet.update(v => !v) },
    { label: 'Theme: Light', key: '', action: () => { themePref.set('light'); localStorage.setItem('mailtap-theme', 'light'); } },
    { label: 'Theme: Dark', key: '', action: () => { themePref.set('dark'); localStorage.setItem('mailtap-theme', 'dark'); } },
    { label: 'Theme: System', key: '', action: () => { themePref.set('system'); localStorage.setItem('mailtap-theme', 'system'); } },
  ];

  $: filtered = query
    ? commands.filter(c => c.label.toLowerCase().includes(query.toLowerCase()))
    : commands;

  $: if (selectedIndex >= filtered.length) selectedIndex = Math.max(0, filtered.length - 1);

  async function handleKeydown(e) {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      selectedIndex = Math.min(selectedIndex + 1, filtered.length - 1);
      await tick();
      scrollSelectedIntoView();
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      selectedIndex = Math.max(selectedIndex - 1, 0);
      await tick();
      scrollSelectedIntoView();
    } else if (e.key === 'Enter') {
      e.preventDefault();
      if (filtered[selectedIndex]) {
        execute(filtered[selectedIndex]);
      }
    }
  }

  function scrollSelectedIntoView() {
    const el = document.querySelector('.command-palette .result.active');
    if (el) el.scrollIntoView({ block: 'nearest' });
  }

  function execute(cmd) {
    showCommandPalette.set(false);
    cmd.action();
  }

  onMount(async () => {
    await tick();
    inputEl?.focus();
  });
</script>

<div class="command-palette-overlay" on:click={() => showCommandPalette.set(false)} on:keydown={handleKeydown}>
  <div class="command-palette" on:click|stopPropagation>
    <input
      bind:this={inputEl}
      bind:value={query}
      placeholder="Type a command..."
      on:keydown={handleKeydown}
    />
    <div class="results">
      {#each filtered as cmd, i}
        <div
          class="result"
          class:active={i === selectedIndex}
          on:click={() => execute(cmd)}
          on:mouseenter={() => selectedIndex = i}
        >
          <span class="label">{cmd.label}</span>
          <span class="hotkey"><kbd>{cmd.key}</kbd></span>
        </div>
      {/each}
      {#if filtered.length === 0}
        <div class="result"><span class="label">No matching commands</span></div>
      {/if}
    </div>
  </div>
</div>
