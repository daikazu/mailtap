<script>
  import { activeTab } from '../lib/stores';

  const tabs = [
    { label: 'HTML Preview', shortcut: '1' },
    { label: 'Plain Text', shortcut: '2' },
    { label: 'Headers', shortcut: '3' },
    { label: 'Attachments', shortcut: '4' },
    { label: 'Source', shortcut: '5' },
  ];
</script>

<div class="tabs">
  {#each tabs as tab, i}
    <button
      class="tab"
      class:active={$activeTab === i}
      on:click={() => activeTab.set(i)}
    >
      {tab.label}<span class="shortcut">{tab.shortcut}</span>
    </button>
  {/each}
</div>
