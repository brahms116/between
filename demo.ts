export interface User {
  age: number;
  $name: string;
  email?: string;
  hobbies?: string[];
  dateOfBirth: string;
  status: Status;
  userData: UserData;
}
export type Status = "Active" | "Disabled" | "pending activation";
export type UserData =
  | { adminData: AdminData }
  | { customerData: CustomerData };
export interface AdminData {
  accessLevel: number;
}
export interface CustomerData {
  attributes: Record<string, unknown>;
}
