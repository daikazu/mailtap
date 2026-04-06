<script>
  import { searchQuery, loadEmails } from '../lib/stores';

  let timer;
  let prevQuery = '';

  function handleInput(e) {
    const value = e.target.value;
    searchQuery.set(value);
    clearTimeout(timer);
    timer = setTimeout(() => {
      loadEmails(value);
    }, 300);
  }

  function handleFocus() {
    prevQuery = $searchQuery;
  }

  function handleKeydown(e) {
    if (e.key === 'Escape') {
      // Revert to previous query and blur
      searchQuery.set(prevQuery);
      clearTimeout(timer);
      loadEmails(prevQuery);
      e.target.blur();
    } else if (e.key === 'Enter') {
      // Accept current filter and blur
      e.target.blur();
    }
  }
</script>

<div style="padding: 8px 12px; border-bottom: 1px solid var(--border); flex-shrink: 0;">
  <input
    class="search"
    placeholder="/ Search emails..."
    value={$searchQuery}
    on:input={handleInput}
    on:focus={handleFocus}
    on:keydown={handleKeydown}
  />
</div>
