import { get } from 'svelte/store';
import { emails, selectedEmailId, selectedEmail, activeTab, showCommandPalette, showHotkeySheet, showConfirmClear, showSettings, selectEmail, deleteEmail, deleteAllEmails } from './stores';

export function initHotkeys() {
  window.addEventListener('keydown', handleKeydown);
  return () => window.removeEventListener('keydown', handleKeydown);
}

function handleKeydown(e) {
  const tag = document.activeElement?.tagName?.toLowerCase();
  const isInput = tag === 'input' || tag === 'textarea';

  // Ctrl+K / Meta+K always works
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault();
    showCommandPalette.update(v => !v);
    return;
  }

  // Esc always works (except when search input is focused — it handles its own Esc)
  if (e.key === 'Escape') {
    if (isInput) return;
    if (get(showSettings)) { showSettings.set(false); return; }
    if (get(showConfirmClear)) { showConfirmClear.set(false); return; }
    if (get(showCommandPalette)) { showCommandPalette.set(false); return; }
    if (get(showHotkeySheet)) { showHotkeySheet.set(false); return; }
    selectEmail(null);
    return;
  }

  // Prevent Tab from cycling through focusable elements (use / to focus search)
  if (e.key === 'Tab') {
    e.preventDefault();
    return;
  }

  // Don't fire other hotkeys when typing in inputs or when a modal is open
  if (isInput) return;
  if (get(showConfirmClear) || get(showCommandPalette) || get(showHotkeySheet) || get(showSettings)) return;

  // / — focus search
  if (e.key === '/') {
    e.preventDefault();
    document.querySelector('.search')?.focus();
    return;
  }

  // j/ArrowDown — next email
  if (e.key === 'j' || e.key === 'ArrowDown') {
    e.preventDefault();
    navigateList(1);
    return;
  }

  // k/ArrowUp — previous email
  if (e.key === 'k' || e.key === 'ArrowUp') {
    e.preventDefault();
    navigateList(-1);
    return;
  }

  // d/Delete/Backspace — delete selected email
  if (e.key === 'd' || e.key === 'Delete' || e.key === 'Backspace') {
    const id = get(selectedEmailId);
    if (id) {
      deleteEmail(id);
    }
    return;
  }

  // h/ArrowLeft — previous tab
  if (e.key === 'h' || e.key === 'ArrowLeft') {
    e.preventDefault();
    activeTab.update(t => Math.max(0, t - 1));
    return;
  }

  // l/ArrowRight — next tab
  if (e.key === 'l' || e.key === 'ArrowRight') {
    e.preventDefault();
    activeTab.update(t => Math.min(4, t + 1));
    return;
  }

  // 1-5 — set active tab
  if (e.key >= '1' && e.key <= '5') {
    activeTab.set(parseInt(e.key) - 1);
    return;
  }

  // Ctrl+Shift+Delete — delete all emails
  if (e.ctrlKey && e.shiftKey && (e.key === 'Delete' || e.key === 'Backspace')) {
    e.preventDefault();
    showConfirmClear.set(true);
    return;
  }

  // , — open settings
  if (e.key === ',') {
    showSettings.set(true);
    return;
  }

  // ? — toggle hotkey sheet
  if (e.key === '?') {
    showHotkeySheet.update(v => !v);
    return;
  }
}

function navigateList(direction) {
  const list = get(emails);
  if (!list || list.length === 0) return;

  const currentId = get(selectedEmailId);
  const currentIndex = list.findIndex(e => e.id === currentId);

  let nextIndex;
  if (currentIndex === -1) {
    nextIndex = direction === 1 ? 0 : list.length - 1;
  } else {
    nextIndex = currentIndex + direction;
    if (nextIndex < 0) nextIndex = 0;
    if (nextIndex >= list.length) nextIndex = list.length - 1;
  }

  selectEmail(list[nextIndex].id);
}
