<script>
  import { emails, selectedEmailId, emailCount, selectEmail, loadEmails, searchQuery, showConfirmClear } from '../lib/stores';
  import SearchBar from './SearchBar.svelte';
  import EmailItem from './EmailItem.svelte';

  function handleClear() {
    showConfirmClear.set(true);
  }

</script>

<div class="list-panel">
  <SearchBar />
  <div class="list-toolbar">
    <button class="btn" on:click={handleClear}>Clear</button>
    <span class="count">{$emailCount} msgs</span>
  </div>
  <div class="email-list">
    {#if $emails.length === 0}
      <div style="padding: 20px; color: var(--text-dim); text-align: center; font-size: 11px;">
        No emails yet.<br/>Send one to localhost:2525
      </div>
    {:else}
      {#each $emails as email (email.id)}
        <EmailItem {email} active={$selectedEmailId === email.id} />
      {/each}
    {/if}
  </div>
</div>
