import { useState } from "react";
import { API_URL } from "./addr";
import type { Commit, Data, Repository } from "./types";

const App = () => {
  const [repo, setRepo] = useState<Repository>([]);
  const [repositoryAddress, setRepositoryAddress] = useState<string>(
    "https://github.com/NarendraPatwardhan/sourcemap",
  );
  const [limit, setLimit] = useState<number>(10);
  const [excludeGlobs, setExcludeGlobs] = useState<string>("");
  const [excludePaths, setExcludePaths] = useState<string>("");

  const fetchData = async () => {
    const FETCH_URL = `${API_URL}/sourcemap`;

    try {
      const body = JSON.stringify({
        address: `${repositoryAddress}.git`,
        limit,
        excludeGlobs,
        excludePaths,
      });
      const response = await fetch(FETCH_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: body,
      });
      const jsonData = await response.json();
      setRepo(jsonData);
    } catch (error) {
      console.error("Error fetching data:", error);
    }
  };

  const renderData = (d: Data[] | null | undefined) => {
    if (!d) return null;
    return (
      <ul>
        {d.map((item: Data, index: number) => (
          <li key={index}>
            <div>
              {item.path}:{item.size}
            </div>
            {renderData(item.children)}
          </li>
        ))}
      </ul>
    );
  };

  return (
    <>
      API URL to be used: {API_URL}
      <hr />
      <form
        onSubmit={(e) => {
          e.preventDefault(); // Prevent the default form submission behavior
          fetchData(); // Call your fetchData function
        }}
      >
        <label>
          Repository Address:
          <input
            type="text"
            placeholder="Enter repository address"
            value={repositoryAddress}
            onChange={(e) => setRepositoryAddress(e.target.value)}
          />
        </label>
        <br />
        <label>
          Integer Field Limit:
          <input
            type="number"
            value={limit}
            onChange={(e) => setLimit(parseInt(e.target.value, 10))}
          />
        </label>
        <br />
        <label>
          Exclude Globs:
          <input
            type="text"
            value={excludeGlobs}
            onChange={(e) => setExcludeGlobs(e.target.value)}
          />
        </label>
        <br />
        <label>
          Exclude Paths:
          <input
            type="text"
            value={excludePaths}
            onChange={(e) => setExcludePaths(e.target.value)}
          />
        </label>
        <br />
        <button type="submit">Fetch data</button>
      </form>
      <hr />
      {repo
        ? (
          <>
            <div>
              {repo.map((entry: Commit, index: number) => {
                return (
                  <div key={index}>
                    {entry.hash.slice(0, 7)}: {entry.message.split("\n")[0]}
                  </div>
                );
              })}
            </div>
            <hr />
            {repo[0]
              ? (
                <div>
                  {renderData(repo[0].data!.children)}
                </div>
              )
              : <p>...</p>}
          </>
        )
        : <p>...</p>}
    </>
  );
};
export default App;
