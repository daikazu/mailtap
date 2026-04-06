<script>
  import { selectedEmail, showHotkeySheet, themePref, cycleTheme } from '../lib/stores';

  const themeIcons = { system: '◑', dark: '●', light: '○' };
  const themeTitles = { system: 'Theme: System', dark: 'Theme: Dark', light: 'Theme: Light' };
</script>

<div class="statusbar">
  {#if $selectedEmail?.attachments?.length > 0}
    <div class="item">{$selectedEmail.attachments.length} attachment(s)</div>
  {/if}
  <div class="spacer"></div>
  <div class="item clickable" on:click={cycleTheme} title={themeTitles[$themePref]}>{themeIcons[$themePref]}</div>
  <div class="item clickable" on:click={() => showHotkeySheet.update(v => !v)}><kbd>?</kbd> help</div>
</div>
