import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { Config, ModelConfig, Message } from '../models/config.model';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = '/api';
  private messageSubject = new BehaviorSubject<Message[]>([]);
  public messages$ = this.messageSubject.asObservable();

  constructor(private http: HttpClient) {}

  // Config APIs
  getConfig(): Observable<Config> {
    return this.http.get<Config>(`${this.apiUrl}/config`);
  }

  saveConfig(config: Config): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/config`, config);
  }

  restartGateway(): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/gateway/restart`, {});
  }

  // Model APIs
  getModels(): Observable<ModelConfig[]> {
    return this.http.get<ModelConfig[]>(`${this.apiUrl}/models`);
  }

  addModel(model: ModelConfig): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/models`, model);
  }

  updateModel(index: number, model: ModelConfig): Observable<void> {
    return this.http.put<void>(`${this.apiUrl}/models/${index}`, model);
  }

  deleteModel(index: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/models/${index}`);
  }

  setDefaultModel(modelName: string): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/models/default`, { model_name: modelName });
  }

  // Chat APIs
  sendMessage(content: string, sessionKey: string = 'web:default'): Observable<any> {
    return this.http.post(`${this.apiUrl}/chat`, {
      content,
      session_key: sessionKey
    });
  }

  // Status
  getStatus(): Observable<any> {
    return this.http.get(`${this.apiUrl}/status`);
  }
}
