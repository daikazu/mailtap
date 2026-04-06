<script>
  import { api } from '../lib/api';

  export let emailId = '';

  let source = '';
  let loading = true;

  $: if (emailId) {
    loading = true;
    api.getRawSource(emailId).then(s => {
      source = s || '';
      loading = false;
    }).catch(() => {
      source = '';
      loading = false;
    });
  }
</script>

{#if loading}
  <div style="color: var(--text-dim); font-size: 12px;">Loading...</div>
{:else}
  <pre style="font-family: var(--font-mono); font-size: 12px; color: var(--text); white-space: pre-wrap; word-wrap: break-word; margin: 0; overflow: auto;">{source || 'No raw source available'}</pre>
{/if}
