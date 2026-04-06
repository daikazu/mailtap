<script>
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import { emails, selectedEmail, selectedEmailId, activeTab, emailCount, searchQuery, showCommandPalette, showHotkeySheet, showConfirmClear, showSettings, loadEmails, selectEmail, smtpPort } from './lib/stores';
  import { initHotkeys } from './lib/hotkeys';
  import EmailList from './components/EmailList.svelte';
  import DetailPanel from './components/DetailPanel.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import CommandPalette from './components/CommandPalette.svelte';
  import HotkeySheet from './components/HotkeySheet.svelte';
  import ConfirmClear from './components/ConfirmClear.svelte';
  import Settings from './components/Settings.svelte';

  onMount(() => {
    const cleanupHotkeys = initHotkeys();
    loadEmails();

    // Listen for new emails from SMTP server (only for external events, not our own actions)
    EventsOn('mailtap:email-received', (summary) => {
      // Only prepend if it matches current search filter (or no filter active)
      const search = get(searchQuery);
      if (!search || summary.subject?.toLowerCase().includes(search.toLowerCase()) ||
          summary.from?.toLowerCase().includes(search.toLowerCase())) {
        emails.update(list => [summary, ...list]);
      }
      emailCount.update(n => n + 1);
      selectEmail(summary.id);
    });

    return cleanupHotkeys;
  });
</script>

<div class="app">
  <div class="titlebar">
    <div class="logo">MAILTAP</div>
    <div class="status">
      <div class="dot"></div>
      SMTP listening on :{$smtpPort}
    </div>
    <div class="cmd"><kbd>Ctrl+K</kbd> Command Palette</div>
  </div>

  <div class="main">
    <EmailList />
    <DetailPanel />
  </div>

  <StatusBar />

  {#if $showCommandPalette}
    <CommandPalette />
  {/if}

  {#if $showHotkeySheet}
    <HotkeySheet />
  {/if}

  {#if $showConfirmClear}
    <ConfirmClear />
  {/if}

  {#if $showSettings}
    <Settings />
  {/if}
</div>
