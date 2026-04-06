<script>
  import { selectEmail } from '../lib/stores';

  export let email;
  export let active = false;

  let el;
  $: if (active && el) {
    el.scrollIntoView({ block: 'nearest' });
  }

  function formatTime(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    const hh = String(d.getHours()).padStart(2, '0');
    const mm = String(d.getMinutes()).padStart(2, '0');
    const ss = String(d.getSeconds()).padStart(2, '0');
    return `${hh}:${mm}:${ss}`;
  }

  function handleClick() {
    selectEmail(email.id);
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
  bind:this={el}
  class="email-item"
  class:active
  class:unread={!email.read}
  on:click={handleClick}
>
  <div class="email-subject">{email.subject || '(no subject)'}</div>
  <div class="email-from">{email.from || ''}</div>
  {#if email.preview}
    <div class="email-preview">{email.preview}</div>
  {/if}
  <div class="email-meta">
    {#if email.attachmentCount > 0}
      <span class="badge attachment">{email.attachmentCount} file{email.attachmentCount > 1 ? 's' : ''}</span>
    {/if}
    <span class="email-time">{formatTime(email.receivedAt)}</span>
  </div>
</div>
