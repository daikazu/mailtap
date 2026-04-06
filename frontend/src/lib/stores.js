import { writable, get, derived } from 'svelte/store';
import { api } from './api';

export const emails = writable([]);
export const selectedEmailId = writable(null);
export const selectedEmail = writable(null);
export const activeTab = writable(0);
export const searchQuery = writable('');
export const emailCount = writable(0);
export const showCommandPalette = writable(false);
export const showHotkeySheet = writable(false);
export const showConfirmClear = writable(false);
export const showSettings = writable(false);
export const smtpPort = writable('2525');

// Theme: 'light', 'dark', or 'system'
const storedThemePref = (typeof localStorage !== 'undefined' && localStorage.getItem('mailtap-theme')) || 'system';
export const themePref = writable(storedThemePref);

// Resolved theme based on pref + system preference
function getSystemTheme() {
  return typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
}

export const resolvedTheme = derived(themePref, ($pref) => {
  return $pref === 'system' ? getSystemTheme() : $pref;
});

export function cycleTheme() {
  const order = ['system', 'dark', 'light'];
  themePref.update(current => {
    const next = order[(order.indexOf(current) + 1) % order.length];
    localStorage.setItem('mailtap-theme', next);
    return next;
  });
}

// Apply theme class to document
if (typeof window !== 'undefined') {
  // Listen for system theme changes
  window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', () => {
    // Force re-derive when system preference changes
    themePref.update(v => v);
  });

  resolvedTheme.subscribe(theme => {
    document.documentElement.setAttribute('data-theme', theme);
  });
}

export async function loadEmails(search = '') {
  const results = await api.listEmails(search, 0, 200);
  emails.set(results || []);
  const count = await api.getEmailCount();
  emailCount.set(count);
  const port = await api.getPort();
  smtpPort.set(port);
}

export async function selectEmail(id) {
  selectedEmailId.set(id);
  if (id) {
    const email = await api.getEmail(id);
    selectedEmail.set(email);
    await api.markAsRead(id);
    emails.update(list => list.map(e => e.id === id ? { ...e, read: true } : e));
  } else {
    selectedEmail.set(null);
  }
}

export async function deleteEmail(id) {
  let nextId = null;
  const list = get(emails);
  const idx = list.findIndex(e => e.id === id);
  if (idx !== -1) {
    if (idx < list.length - 1) nextId = list[idx + 1].id;
    else if (idx > 0) nextId = list[idx - 1].id;
  }

  await api.deleteEmail(id);
  emails.update(l => l.filter(e => e.id !== id));
  const count = await api.getEmailCount();
  emailCount.set(count);

  if (nextId) {
    await selectEmail(nextId);
  } else {
    selectedEmailId.set(null);
    selectedEmail.set(null);
  }
}

export async function deleteAllEmails() {
  await api.deleteAllEmails();
  emails.set([]);
  selectedEmailId.set(null);
  selectedEmail.set(null);
  emailCount.set(0);
}
