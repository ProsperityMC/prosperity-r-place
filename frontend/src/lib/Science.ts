export interface Science {
  names: Map<string, string>;
  users: Map<string, ScienceUser>;
}

export interface ScienceUser {
  shape: string;
  selArea: string;
}
