function GetUniqueTotalPage(arr: any, index: any) {
    const grouped = arr
        .reduce((accumulator: any, currentItem: any) => {
        if (!accumulator[currentItem[index]]) {
            accumulator[currentItem[index]] = [];
        }
        accumulator[currentItem[index]].push(currentItem);
        return accumulator;
        }, {});

    let result = Object.values(grouped).map((group: any) => group.reduce((maxTotalPagesItem: any, currentItem: any) => {
        return parseInt(currentItem.TotalPages) > parseInt(maxTotalPagesItem.TotalPages) ? currentItem : maxTotalPagesItem;
    }));

    return result;
}
  

export default GetUniqueTotalPage