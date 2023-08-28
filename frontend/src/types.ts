interface Commit {
  hash: string;
  author: string;
  message: string;
  timestamp: string;
  data: Data | null;
}

interface Data {
  name: string;
  path: string;
  size: string;
  changes: Changes;
  content?: string; // Optional property
  repr?: string; // Optional property
  children?: Data[] | null;
}

interface Changes {
  addition: number;
  deletion: number;
}

type Repository = Commit[];

export type { Changes, Commit, Data, Repository };
