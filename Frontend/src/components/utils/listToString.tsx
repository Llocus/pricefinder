const ListToString = (list: any) => {
    let text = "";
    // eslint-disable-next-line array-callback-return
    list.map((l: any) => {
        text += "," + l.toString().toLowerCase()
    })
    text = text.substring(1);
    return text
}

export default ListToString;