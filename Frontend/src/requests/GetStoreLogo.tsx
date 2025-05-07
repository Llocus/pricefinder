const GetStoreLogo = ({ props }: any) => {
    return fetch(process.env.REACT_APP_SERVER_URL + props.logo)
    .then(response => response.json())
    .then((data) => data)
    .catch((error: Error) => console.error(error));
}

export default GetStoreLogo