import axios from "axios";
import Store from "../types/Store";

const GetStores = () => {
    return axios.get(process.env.REACT_APP_SERVER_URL + "/stores")
        .then(response => response.data)
        .then((data: Store[]) => data)
        .catch((error: Error) => console.error(error));
}

export default GetStores