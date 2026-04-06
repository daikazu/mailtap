<script>
  import { selectedEmail, activeTab } from '../lib/stores';
  import { deleteEmail } from '../lib/stores';
  import DetailHeader from './DetailHeader.svelte';
  import TabBar from './TabBar.svelte';
  import HtmlPreview from './HtmlPreview.svelte';
  import PlainText from './PlainText.svelte';
  import Headers from './Headers.svelte';
  import Attachments from './Attachments.svelte';
  import RawSource from './RawSource.svelte';
</script>

<div class="detail-panel">
  {#if $selectedEmail}
    <DetailHeader email={$selectedEmail} on:delete={() => deleteEmail($selectedEmail.id)} />
    <TabBar />
    <div class="preview">
      {#if $activeTab === 0}
        <HtmlPreview html={$selectedEmail.htmlBody} />
      {:else if $activeTab === 1}
        <PlainText text={$selectedEmail.plainBody} />
      {:else if $activeTab === 2}
        <Headers headers={$selectedEmail.headers} />
      {:else if $activeTab === 3}
        <Attachments attachments={$selectedEmail.attachments || []} />
      {:else if $activeTab === 4}
        <RawSource emailId={$selectedEmail.id} />
      {/if}
    </div>
  {:else}
    <div style="display: flex; align-items: center; justify-content: center; height: 100%; color: var(--text-dim); font-size: 12px;">
      Select an email to preview
    </div>
  {/if}
</div>
