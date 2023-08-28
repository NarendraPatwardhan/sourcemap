import { useState } from "react";
import { API_URL } from "./addr";
import type { Commit, Repository } from "./types";

const App = () => {
  const [data, setData] = useState<Repository>([]);
  const [repositoryAddress, setRepositoryAddress] = useState(
    "https://github.com/NarendraPatwardhan/sourcemap",
  );

  const fetchData = async () => {
    const FETCH_URL = `${API_URL}/sourcemap`;

    try {
      const body = JSON.stringify({ address: `${repositoryAddress}.git` });
      const response = await fetch(FETCH_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });
      const jsonData = await response.json();
      setData(jsonData);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  return (
    <>
      API URL to be used: {API_URL}
      <hr />
      <input
        type="text"
        placeholder="Enter repository address"
        value={repositoryAddress}
        onChange={(e) => setRepositoryAddress(e.target.value)}
      />
      <button onClick={fetchData}>Fetch data</button>
      <hr />
      {data
        ? (
          <div>
            {data.map((entry: Commit, index: number) => {
              return (
                <div key={index}>
                  {entry.hash.slice(0, 7)}:
                  {entry.message.split("\n")[0]}
                </div>
              );
            })}
          </div>
        )
        : <p>...</p>}
    </>
  );
};

export default App;
