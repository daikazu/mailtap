package notify

/*
#cgo CFLAGS: -x objective-c -mmacosx-version-min=10.14
#cgo LDFLAGS: -framework Foundation -framework UserNotifications
#import <UserNotifications/UserNotifications.h>

static void MailTapSendNotification(const char *title, const char *body) {
    @autoreleasepool {
        if (@available(macOS 10.14, *)) {
            UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];

            [center requestAuthorizationWithOptions:(UNAuthorizationOptionAlert | UNAuthorizationOptionSound)
                completionHandler:^(BOOL granted, NSError *error) {
                    if (!granted) return;

                    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
                    content.title = [NSString stringWithUTF8String:title];
                    content.body = [NSString stringWithUTF8String:body];
                    content.sound = [UNNotificationSound defaultSound];

                    UNNotificationRequest *request = [UNNotificationRequest
                        requestWithIdentifier:[[NSUUID UUID] UUIDString]
                        content:content
                        trigger:nil];

                    [center addNotificationRequest:request withCompletionHandler:nil];
                }];
        }
    }
}
*/
import "C"

func Send(title, body string) {
	C.MailTapSendNotification(C.CString(title), C.CString(body))
}
