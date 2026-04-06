<script>
  import { api } from '../lib/api';
  export let attachments = [];

  function formatSize(bytes) {
    if (!bytes) return '0 B';
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }

  async function save(id) {
    try {
      await api.saveAttachment(id);
    } catch (e) {
      console.error('Failed to save attachment:', e);
    }
  }
</script>

<div style="font-size: 12px;">
  {#if attachments.length === 0}
    <p style="color: var(--text-dim);">No attachments</p>
  {:else}
    {#each attachments as att}
      <div style="display: flex; align-items: center; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid var(--border);">
        <div>
          <div style="color: var(--text-bright); font-weight: 500;">{att.filename || 'Untitled'}</div>
          <div style="color: var(--text-dim); font-size: 11px; margin-top: 2px;">
            {att.contentType || 'unknown'} — {formatSize(att.size)}
          </div>
        </div>
        <button
          on:click={() => save(att.id)}
          style="background: none; border: 1px solid var(--border); color: var(--text-dim); padding: 4px 10px; border-radius: 4px; cursor: pointer; font-size: 11px; font-family: var(--font-mono);"
        >
          Save
        </button>
      </div>
    {/each}
  {/if}
</div>
