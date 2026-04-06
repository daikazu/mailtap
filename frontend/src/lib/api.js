import { ListEmails, GetEmail, DeleteEmail, DeleteAllEmails,
         GetRawSource, MarkAsRead, GetEmailCount, SaveAttachment,
         GetPort, GetSettings, SaveSettings
} from '../../wailsjs/go/main/App';

export const api = {
  listEmails: (search = '', offset = 0, limit = 50) => ListEmails(search, offset, limit),
  getEmail: (id) => GetEmail(id),
  deleteEmail: (id) => DeleteEmail(id),
  deleteAllEmails: () => DeleteAllEmails(),
  getRawSource: (id) => GetRawSource(id),
  markAsRead: (id) => MarkAsRead(id),
  getEmailCount: () => GetEmailCount(),
  saveAttachment: (attachmentId) => SaveAttachment(attachmentId),
  getPort: () => GetPort(),
  getSettings: () => GetSettings(),
  saveSettings: (s) => SaveSettings(s),
};
