interface User {
  age: number;
  name: string;
  email?: string;
  hobbies?: string[];
  status: Status;
  userData: UserData;
}
type Status = "Active" | "Disabled" | "Pending";
type UserData =
  | { _type: "adminData"; adminData: AdminData }
  | { _type: "customerData"; customerData: CustomerData };
interface AdminData {
  accessLevel: number;
}
interface CustomerData {
  attributes: Record<string, unknown>;
}
