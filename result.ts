interface User {
  name: string;
  age?: number;
  gender?: string;
  friendIds?: (string | undefined)[];
  array2d: ((number | undefined)[] | undefined)[];
  status: Status;
  data: Data;
}

type Status = "NOT_HERE" | "REALLY";

type Data =
  | { _type: "adminData"; adminData: AdminData }
  | { _type: "userData"; userData: UserData };

interface AdminData {
  level: number;
}
interface UserData {
  id: number;
  data: Record<string, unknown>;
}
