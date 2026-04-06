<script>
  import { showConfirmClear, deleteAllEmails } from '../lib/stores';

  function handleConfirm() {
    showConfirmClear.set(false);
    deleteAllEmails();
  }

  function handleCancel() {
    showConfirmClear.set(false);
  }

  function handleKeydown(e) {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleConfirm();
    }
  }
</script>

<div class="command-palette-overlay" on:click={handleCancel} on:keydown={handleKeydown}>
  <div class="command-palette" on:click|stopPropagation style="width: 360px;">
    <div style="padding: 20px 24px;">
      <div style="font-size: 14px; color: var(--text-bright); margin-bottom: 4px;">Delete all emails?</div>
      <div style="font-size: 12px; color: var(--text-dim); margin-bottom: 20px;">This cannot be undone.</div>
      <div style="display: flex; gap: 8px; justify-content: flex-end;">
        <button class="btn" on:click={handleCancel}>Cancel <kbd>Esc</kbd></button>
        <button class="btn danger" on:click={handleConfirm}>Delete All <kbd>Enter</kbd></button>
      </div>
    </div>
  </div>
</div>
