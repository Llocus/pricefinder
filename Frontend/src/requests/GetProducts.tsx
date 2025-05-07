import axios, { AxiosRequestConfig } from "axios";
import Product from "../types/Product";

const GetProducts = ({ params }: AxiosRequestConfig["params"]) => {
        return axios.get(process.env.REACT_APP_SERVER_URL + "/", { params, withCredentials: false, headers: { Accept: true } })
                .then(response => response.data)
                .then((data: Product[]) => data)
                .catch((error: Error) => console.error(error));
}

export default GetProducts;