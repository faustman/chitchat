import { AuthService } from "../Auth/service";

export class ChatService {
  static fetchMessages(): Promise<Array<any>> {
    return fetch("/messages?token=" + AuthService.token)
      .then((r) => r.json()).then((data) => data.messages || []);
  }

  static fetchUsers(): Promise<Array<any>> {
    return fetch("/users?token=" + AuthService.token).then((r) => r.json()).then((data) => data.users || []);
  }
}