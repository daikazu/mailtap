<script>
  import { createEventDispatcher } from 'svelte';
  export let email;
  const dispatch = createEventDispatcher();

  function formatDate(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    return d.toLocaleString();
  }

  function formatSize(bytes) {
    if (!bytes) return '';
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }
</script>

<div class="detail-header">
  <div style="display: flex; justify-content: space-between; align-items: flex-start;">
    <div style="flex: 1; min-width: 0;">
      <div class="subject">{email.subject || '(no subject)'}</div>
      <div class="meta">
        <div>
          <span>From:</span> {email.from}
          <span style="margin-left: 8px;">To:</span> {email.to}
          {#if email.cc}
            <span style="margin-left: 8px;">CC:</span> {email.cc}
          {/if}
        </div>
        <div>
          {#if email.date}<span>Date:</span> {formatDate(email.date)}{/if}
          {#if email.size}<span style="margin-left: 12px;">Size:</span> {formatSize(email.size)}{/if}
          {#if email.messageId}<span style="margin-left: 12px;">ID:</span> {email.messageId}{/if}
        </div>
      </div>
    </div>
    <button
      on:click={() => dispatch('delete')}
      style="background: none; border: 1px solid var(--border); color: var(--text-dim); padding: 4px 10px; border-radius: 4px; cursor: pointer; font-size: 11px; font-family: var(--font-mono); flex-shrink: 0; margin-left: 12px;"
    >
      Delete
    </button>
  </div>
</div>
