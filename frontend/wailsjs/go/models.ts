export namespace models {
	
	export class ActivityLog {
	    id: string;
	    user_id: string;
	    user_name?: string;
	    action: string;
	    entity_type: string;
	    entity_id?: string;
	    details?: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new ActivityLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.action = source["action"];
	        this.entity_type = source["entity_type"];
	        this.entity_id = source["entity_id"];
	        this.details = source["details"];
	        this.created_at = this.convertValues(source["created_at"], null);
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
	export class ActivityListResponse {
	    success: boolean;
	    message: string;
	    activities?: ActivityLog[];
	
	    static createFrom(source: any = {}) {
	        return new ActivityListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.activities = this.convertValues(source["activities"], ActivityLog);
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
	
	export class TaskAttachment {
	    id: string;
	    task_id: string;
	    user_id: string;
	    user_name?: string;
	    file_name: string;
	    file_path: string;
	    file_size: number;
	    // Go type: time
	    uploaded_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TaskAttachment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task_id = source["task_id"];
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.file_name = source["file_name"];
	        this.file_path = source["file_path"];
	        this.file_size = source["file_size"];
	        this.uploaded_at = this.convertValues(source["uploaded_at"], null);
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
	export class AttachmentListResponse {
	    success: boolean;
	    message: string;
	    attachments?: TaskAttachment[];
	
	    static createFrom(source: any = {}) {
	        return new AttachmentListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.attachments = this.convertValues(source["attachments"], TaskAttachment);
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
	export class AttachmentResponse {
	    success: boolean;
	    message: string;
	    attachment?: TaskAttachment;
	
	    static createFrom(source: any = {}) {
	        return new AttachmentResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.attachment = this.convertValues(source["attachment"], TaskAttachment);
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
	export class User {
	    id: string;
	    email: string;
	    full_name: string;
	    avatar_url?: string;
	    timezone: string;
	    notification_preferences: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.full_name = source["full_name"];
	        this.avatar_url = source["avatar_url"];
	        this.timezone = source["timezone"];
	        this.notification_preferences = source["notification_preferences"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class AuthState {
	    is_logged_in: boolean;
	    user?: User;
	
	    static createFrom(source: any = {}) {
	        return new AuthState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_logged_in = source["is_logged_in"];
	        this.user = this.convertValues(source["user"], User);
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
	export class BaseResponse {
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new BaseResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}
	export class ChangePasswordRequest {
	    old_password: string;
	    new_password: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangePasswordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.old_password = source["old_password"];
	        this.new_password = source["new_password"];
	    }
	}
	export class TaskComment {
	    id: string;
	    task_id: string;
	    user_id: string;
	    user_name?: string;
	    content: string;
	    reactions?: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TaskComment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task_id = source["task_id"];
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.content = source["content"];
	        this.reactions = source["reactions"];
	        this.created_at = this.convertValues(source["created_at"], null);
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
	export class CommentListResponse {
	    success: boolean;
	    message: string;
	    comments?: TaskComment[];
	
	    static createFrom(source: any = {}) {
	        return new CommentListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.comments = this.convertValues(source["comments"], TaskComment);
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
	export class CommentResponse {
	    success: boolean;
	    message: string;
	    comment?: TaskComment;
	
	    static createFrom(source: any = {}) {
	        return new CommentResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.comment = this.convertValues(source["comment"], TaskComment);
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
	export class CreateCommentRequest {
	    task_id: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateCommentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task_id = source["task_id"];
	        this.content = source["content"];
	    }
	}
	export class CreateProjectRequest {
	    name: string;
	    description: string;
	    color: string;
	    template?: string;
	    start_date?: string;
	    end_date?: string;
	    total_budget?: number;
	    currency?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.color = source["color"];
	        this.template = source["template"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.total_budget = source["total_budget"];
	        this.currency = source["currency"];
	    }
	}
	export class CreateTaskRequest {
	    project_id: string;
	    parent_id?: string;
	    name: string;
	    description: string;
	    start_date?: string;
	    end_date?: string;
	    duration: number;
	    estimated_hours: number;
	    status: string;
	    priority: string;
	    labels?: string;
	    checklists?: string;
	    estimated_cost?: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.parent_id = source["parent_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.duration = source["duration"];
	        this.estimated_hours = source["estimated_hours"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.labels = source["labels"];
	        this.checklists = source["checklists"];
	        this.estimated_cost = source["estimated_cost"];
	    }
	}
	export class CreateTimeLogRequest {
	    task_id: string;
	    logged_date: string;
	    hours: number;
	    description?: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateTimeLogRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.task_id = source["task_id"];
	        this.logged_date = source["logged_date"];
	        this.hours = source["hours"];
	        this.description = source["description"];
	    }
	}
	export class LoginRequest {
	    email: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.email = source["email"];
	        this.password = source["password"];
	    }
	}
	export class LoginResponse {
	    success: boolean;
	    message: string;
	    user?: User;
	    token?: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.user = this.convertValues(source["user"], User);
	        this.token = source["token"];
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
	export class MemberWorkload {
	    user_id: string;
	    user_name: string;
	    date: string;
	    allocated_hours: number;
	    daily_capacity_hours: number;
	
	    static createFrom(source: any = {}) {
	        return new MemberWorkload(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.date = source["date"];
	        this.allocated_hours = source["allocated_hours"];
	        this.daily_capacity_hours = source["daily_capacity_hours"];
	    }
	}
	export class Notification {
	    id: string;
	    user_id: string;
	    title: string;
	    message: string;
	    is_read: boolean;
	    type: string;
	    related_id: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Notification(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_id = source["user_id"];
	        this.title = source["title"];
	        this.message = source["message"];
	        this.is_read = source["is_read"];
	        this.type = source["type"];
	        this.related_id = source["related_id"];
	        this.created_at = this.convertValues(source["created_at"], null);
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
	export class NotificationListResponse {
	    success: boolean;
	    message: string;
	    notifications?: Notification[];
	    unread_count: number;
	
	    static createFrom(source: any = {}) {
	        return new NotificationListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.notifications = this.convertValues(source["notifications"], Notification);
	        this.unread_count = source["unread_count"];
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
	export class NotificationResponse {
	    success: boolean;
	    message: string;
	    notification?: Notification;
	
	    static createFrom(source: any = {}) {
	        return new NotificationResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.notification = this.convertValues(source["notification"], Notification);
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
	export class Project {
	    id: string;
	    workspace_id: string;
	    name: string;
	    description: string;
	    color: string;
	    status: string;
	    // Go type: time
	    start_date?: any;
	    // Go type: time
	    end_date?: any;
	    is_public: boolean;
	    share_token?: string;
	    total_budget: number;
	    currency: string;
	    created_by: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.workspace_id = source["workspace_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.color = source["color"];
	        this.status = source["status"];
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.end_date = this.convertValues(source["end_date"], null);
	        this.is_public = source["is_public"];
	        this.share_token = source["share_token"];
	        this.total_budget = source["total_budget"];
	        this.currency = source["currency"];
	        this.created_by = source["created_by"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class ProjectListResponse {
	    success: boolean;
	    message: string;
	    projects?: Project[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.projects = this.convertValues(source["projects"], Project);
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
	export class ProjectMember {
	    project_id: string;
	    user_id: string;
	    role: string;
	    daily_capacity_hours: number;
	    user_email?: string;
	    user_full_name?: string;
	    user_avatar?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectMember(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.user_id = source["user_id"];
	        this.role = source["role"];
	        this.daily_capacity_hours = source["daily_capacity_hours"];
	        this.user_email = source["user_email"];
	        this.user_full_name = source["user_full_name"];
	        this.user_avatar = source["user_avatar"];
	    }
	}
	export class ProjectMemberListResponse {
	    success: boolean;
	    message: string;
	    members?: ProjectMember[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectMemberListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.members = this.convertValues(source["members"], ProjectMember);
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
	export class ProjectMemberRequest {
	    project_id: string;
	    user_id: string;
	    role: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectMemberRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.user_id = source["user_id"];
	        this.role = source["role"];
	    }
	}
	export class ProjectMemberResponse {
	    success: boolean;
	    message: string;
	    member?: ProjectMember;
	
	    static createFrom(source: any = {}) {
	        return new ProjectMemberResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.member = this.convertValues(source["member"], ProjectMember);
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
	export class ProjectResponse {
	    success: boolean;
	    message: string;
	    project?: Project;
	
	    static createFrom(source: any = {}) {
	        return new ProjectResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.project = this.convertValues(source["project"], Project);
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
	export class RegisterRequest {
	    full_name: string;
	    email: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new RegisterRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full_name = source["full_name"];
	        this.email = source["email"];
	        this.password = source["password"];
	    }
	}
	export class Task {
	    id: string;
	    project_id: string;
	    parent_id?: string;
	    name: string;
	    description: string;
	    status: string;
	    priority: string;
	    // Go type: time
	    start_date?: any;
	    // Go type: time
	    end_date?: any;
	    duration: number;
	    estimated_hours: number;
	    progress: number;
	    assignee_id?: string;
	    sort_order: number;
	    dependencies: string;
	    is_milestone: boolean;
	    labels?: string;
	    checklists?: string;
	    estimated_cost: number;
	    actual_cost: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.parent_id = source["parent_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.start_date = this.convertValues(source["start_date"], null);
	        this.end_date = this.convertValues(source["end_date"], null);
	        this.duration = source["duration"];
	        this.estimated_hours = source["estimated_hours"];
	        this.progress = source["progress"];
	        this.assignee_id = source["assignee_id"];
	        this.sort_order = source["sort_order"];
	        this.dependencies = source["dependencies"];
	        this.is_milestone = source["is_milestone"];
	        this.labels = source["labels"];
	        this.checklists = source["checklists"];
	        this.estimated_cost = source["estimated_cost"];
	        this.actual_cost = source["actual_cost"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	
	
	export class TaskListResponse {
	    success: boolean;
	    message: string;
	    tasks?: Task[];
	
	    static createFrom(source: any = {}) {
	        return new TaskListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.tasks = this.convertValues(source["tasks"], Task);
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
	export class TaskResponse {
	    success: boolean;
	    message: string;
	    task?: Task;
	
	    static createFrom(source: any = {}) {
	        return new TaskResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.task = this.convertValues(source["task"], Task);
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
	export class TimeLog {
	    id: string;
	    task_id: string;
	    user_id: string;
	    user_name?: string;
	    logged_date: string;
	    hours: number;
	    description?: string;
	    // Go type: time
	    created_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TimeLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task_id = source["task_id"];
	        this.user_id = source["user_id"];
	        this.user_name = source["user_name"];
	        this.logged_date = source["logged_date"];
	        this.hours = source["hours"];
	        this.description = source["description"];
	        this.created_at = this.convertValues(source["created_at"], null);
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
	export class TimeLogListResponse {
	    success: boolean;
	    message: string;
	    time_logs?: TimeLog[];
	
	    static createFrom(source: any = {}) {
	        return new TimeLogListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.time_logs = this.convertValues(source["time_logs"], TimeLog);
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
	export class TimeLogResponse {
	    success: boolean;
	    message: string;
	    time_log?: TimeLog;
	
	    static createFrom(source: any = {}) {
	        return new TimeLogResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.time_log = this.convertValues(source["time_log"], TimeLog);
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
	export class UpdateCapacityRequest {
	    project_id: string;
	    user_id: string;
	    daily_capacity_hours: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCapacityRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.project_id = source["project_id"];
	        this.user_id = source["user_id"];
	        this.daily_capacity_hours = source["daily_capacity_hours"];
	    }
	}
	export class UpdateCommentRequest {
	    id: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCommentRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	    }
	}
	export class UpdateProfileRequest {
	    full_name: string;
	    avatar_url?: string;
	    timezone: string;
	    notification_preferences?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProfileRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full_name = source["full_name"];
	        this.avatar_url = source["avatar_url"];
	        this.timezone = source["timezone"];
	        this.notification_preferences = source["notification_preferences"];
	    }
	}
	export class UpdateProjectRequest {
	    id: string;
	    name: string;
	    description: string;
	    color: string;
	    status: string;
	    start_date?: string;
	    end_date?: string;
	    total_budget?: number;
	    currency?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProjectRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.color = source["color"];
	        this.status = source["status"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.total_budget = source["total_budget"];
	        this.currency = source["currency"];
	    }
	}
	export class UpdateTaskRequest {
	    id: string;
	    name: string;
	    status: string;
	    priority: string;
	    duration: number;
	    start_date?: string;
	    end_date?: string;
	    assignee_id?: string;
	    description?: string;
	    progress?: number;
	    sort_order?: number;
	    estimated_hours?: number;
	    dependencies?: string;
	    is_milestone?: boolean;
	    labels?: string;
	    checklists?: string;
	    estimated_cost?: number;
	    actual_cost?: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	        this.duration = source["duration"];
	        this.start_date = source["start_date"];
	        this.end_date = source["end_date"];
	        this.assignee_id = source["assignee_id"];
	        this.description = source["description"];
	        this.progress = source["progress"];
	        this.sort_order = source["sort_order"];
	        this.estimated_hours = source["estimated_hours"];
	        this.dependencies = source["dependencies"];
	        this.is_milestone = source["is_milestone"];
	        this.labels = source["labels"];
	        this.checklists = source["checklists"];
	        this.estimated_cost = source["estimated_cost"];
	        this.actual_cost = source["actual_cost"];
	    }
	}
	export class UpdateTimeLogRequest {
	    id: string;
	    logged_date: string;
	    hours: number;
	    description?: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateTimeLogRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.logged_date = source["logged_date"];
	        this.hours = source["hours"];
	        this.description = source["description"];
	    }
	}
	
	export class UserListResponse {
	    success: boolean;
	    message: string;
	    users?: User[];
	
	    static createFrom(source: any = {}) {
	        return new UserListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.users = this.convertValues(source["users"], User);
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
	export class WorkloadReportResponse {
	    success: boolean;
	    message: string;
	    workloads?: MemberWorkload[];
	
	    static createFrom(source: any = {}) {
	        return new WorkloadReportResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.workloads = this.convertValues(source["workloads"], MemberWorkload);
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

