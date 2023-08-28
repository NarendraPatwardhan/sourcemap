import { useEffect, useState } from "react";
import { API_URL } from "./addr";
import type { Commit, Repository } from "./types";
import "./App.css";

const App = () => {
  const [data, setData] = useState<Repository>([]);

  useEffect(() => {
    const fetchData = async () => {
      const FETCH_URL = `${API_URL}/sourcemap`;

      try {
        const response = await fetch(FETCH_URL);
        const jsonData = await response.json();
        setData(jsonData);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };
    fetchData();
  }, []);

  return (
    <>
      API URL to be used: {API_URL}
      <hr />
      {data
        ? (
          <div>
            {data.map((entry: Commit, index: number) => {
              return (
                <div key={index}>
                  {entry.hash.slice(0, 7)} :
                  {entry.message.split("\n")[0]}
                </div>
              );
            })}
          </div>
        )
        : <p>Loading data...</p>}
    </>
  );
};

export default App;
