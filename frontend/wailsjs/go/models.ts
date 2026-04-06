export namespace models {
	
	export class Attachment {
	    id: string;
	    emailId: string;
	    filename: string;
	    contentType: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new Attachment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.emailId = source["emailId"];
	        this.filename = source["filename"];
	        this.contentType = source["contentType"];
	        this.size = source["size"];
	    }
	}
	export class Email {
	    id: string;
	    from: string;
	    to: string[];
	    cc: string[];
	    bcc: string[];
	    subject: string;
	    htmlBody: string;
	    plainBody: string;
	    headers: Record<string, Array<string>>;
	    attachments: Attachment[];
	    // Go type: time
	    receivedAt: any;
	    read: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Email(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.from = source["from"];
	        this.to = source["to"];
	        this.cc = source["cc"];
	        this.bcc = source["bcc"];
	        this.subject = source["subject"];
	        this.htmlBody = source["htmlBody"];
	        this.plainBody = source["plainBody"];
	        this.headers = source["headers"];
	        this.attachments = this.convertValues(source["attachments"], Attachment);
	        this.receivedAt = this.convertValues(source["receivedAt"], null);
	        this.read = source["read"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EmailSummary {
	    id: string;
	    from: string;
	    to: string[];
	    subject: string;
	    preview: string;
	    attachmentCount: number;
	    // Go type: time
	    receivedAt: any;
	    read: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EmailSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.from = source["from"];
	        this.to = source["to"];
	        this.subject = source["subject"];
	        this.preview = source["preview"];
	        this.attachmentCount = source["attachmentCount"];
	        this.receivedAt = this.convertValues(source["receivedAt"], null);
	        this.read = source["read"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace settings {
	
	export class Settings {
	    port: string;
	    notifications: boolean;
	    theme: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.notifications = source["notifications"];
	        this.theme = source["theme"];
	    }
	}

}

