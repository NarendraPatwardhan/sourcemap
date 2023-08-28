import { useEffect, useState } from "react";
import { API_URL } from "./addr";
import "./App.css";

function App() {
  const [data, setData] = useState(null);

  useEffect(() => {
    const FETCH_URL = `${API_URL}/sourcemap`;

    fetch(FETCH_URL)
      .then((response) => response.json())
      .then((jsonData) => setData(jsonData))
      .catch((error) => console.error("Error fetching data:", error));
  }, []);

  return (
    <>
      API URL to be used : {API_URL}
      {data
        ? <pre>{JSON.stringify(data[0], null, 2)}</pre>
        : <p>Loading data...</p>}
    </>
  );
}

export default App;
