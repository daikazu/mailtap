<script>
  export let headers = {};
</script>

<div style="font-family: var(--font-mono); font-size: 12px; overflow: auto;">
  {#if headers && Object.keys(headers).length > 0}
    <table style="width: 100%; border-collapse: collapse;">
      {#each Object.entries(headers) as [key, values]}
        <tr>
          <td style="color: var(--accent); padding: 4px 12px 4px 0; vertical-align: top; white-space: nowrap; font-weight: 600;">{key}</td>
          <td style="color: var(--text); padding: 4px 0; word-break: break-all;">{Array.isArray(values) ? values.join(', ') : values}</td>
        </tr>
      {/each}
    </table>
  {:else}
    <p style="color: var(--text-dim);">No headers available</p>
  {/if}
</div>
