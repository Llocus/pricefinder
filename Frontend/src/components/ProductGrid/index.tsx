import { Box, Center, SimpleGrid, Spinner, Text } from "@chakra-ui/react";
import ProductCard from "../ProductCard";
import Store from "../../types/Store";
import { useEffect, useState } from "react";
import Product from "../../types/Product";
import GetProducts from "../../requests/GetProducts";
import ListToString from "../utils/listToString";
import useWindowDimensions from "../../hooks/useWindowDimensions";
import md5 from "md5";
import GetUniqueTotalPage from "../utils/getUniqueMaxPage";

const ProductGrid = ({
    isLoadingStores,
    isLoadingProducts,
    setIsLoadingProducts,
    isRequestingNextPage,
    setIsRequestingNextPage,
    stores,
    products,
    setProducts,
    search,
    minPrice,
    maxPrice,
    storesFilter,
    updateFilters,
    setUpdateFilters,
    forceUpdateFilters,
    setForceUpdateFilters,
    gridRef }:
    {
        isLoadingStores: boolean,
        isLoadingProducts: boolean,
        setIsLoadingProducts: React.Dispatch<React.SetStateAction<boolean>>,
        isRequestingNextPage: boolean,
        setIsRequestingNextPage: React.Dispatch<React.SetStateAction<boolean>>,
        stores: Store[] | undefined,
        products: Product[] | undefined,
        setProducts: React.Dispatch<React.SetStateAction<Product[] | undefined>>,
        search: string,
        gridRef: React.MutableRefObject<null>,
        minPrice: string,
        maxPrice: string,
        updateFilters: boolean,
        setUpdateFilters: React.Dispatch<React.SetStateAction<boolean>>,
        forceUpdateFilters: boolean,
        setForceUpdateFilters: React.Dispatch<React.SetStateAction<boolean>>,
        storesFilter: string[]
    }) => {
    const [page, setPage] = useState(1)
    const [lastItemCode, setLastItemCode] = useState("card-0")
    const [lastFilter, setLastFilter] = useState(md5(`${minPrice}${maxPrice}${ListToString(storesFilter)}`));

    const { height } = useWindowDimensions()

    function removeDuplicates(array: Product[]) {
        return array.filter((obj: any, index: any) => {
            const objString = JSON.stringify(obj);

            return (
                index ===
                array.findIndex((objB: Product) => {
                    return JSON.stringify(objB) === objString;
                })
            );
        });
    }

    const handleNextPage = () => {
        if (products) {
            let nextPageStoresList = [""];
            GetUniqueTotalPage(products, "store").map((p: Product) => {
                if (parseInt(p.totalPages) >= page) {
                    nextPageStoresList.push(p.store)
                }
                return p
            })
            if (ListToString(nextPageStoresList.filter(n => n).filter((ns) => !storesFilter.includes(ns))).length) {
                GetProducts({
                    params: {
                        "pg": page,
                        "q": search,
                        "pi": minPrice,
                        "pf": maxPrice,
                        "stores": ListToString(nextPageStoresList.filter(n => n).filter((ns) => !storesFilter.includes(ns))?.map((s) => {
                            return s
                        })),
                    }
                }).then((r) => {
                    if (r) {
                        setPage(page + 1)
                        setProducts(removeDuplicates([...products, ...r]))
                        const code = document.querySelector(lastItemCode)
                        if (code) {
                            code.scrollIntoView({ block: "end", behavior: "smooth" })
                        }
                        setIsRequestingNextPage(false)
                    }
                })
            }
        }
    }

    useEffect(() => {
        const actualFilter = md5(`${minPrice}${maxPrice}${ListToString(storesFilter)}`)
        if ((updateFilters && lastFilter !== actualFilter) || forceUpdateFilters) {
            setLastFilter(actualFilter)
            setIsLoadingProducts(true)
            setPage(1)
            setProducts(undefined)
            GetProducts({
                params: {
                    "pg": page,
                    "q": search,
                    "pi": minPrice ?? "",
                    "pf": maxPrice ?? "",
                    "stores": ListToString(stores?.filter((s) => !storesFilter.includes(s.name))?.map((s) => {
                        return s.name
                    })),
                }
            }).then((r) => {
                if (r) {
                    setProducts(r)
                    setIsLoadingProducts(false)
                    setUpdateFilters(false)
                    setForceUpdateFilters(false)
                }
            })
        } else {
            setUpdateFilters(false)
        }
    }, [forceUpdateFilters, lastFilter, maxPrice, minPrice, page, search, setForceUpdateFilters, setIsLoadingProducts, setProducts, setUpdateFilters, stores, storesFilter, updateFilters])

    const handleScroll = (e: React.UIEvent<HTMLElement>) => {
        const target = e.currentTarget;
        let nextPageStoresList = [""];
        GetUniqueTotalPage(products, "store").map((p: Product) => {
            if (parseInt(p.totalPages) >= page) {
                nextPageStoresList.push(p.store)
            }
            return p
        })
        if (products) {
            setLastItemCode(`.card-${products.length}`)
        }
        var maxScrollLeft = target?.scrollHeight - target?.clientHeight;
        if (target?.scrollTop > 100 && maxScrollLeft > 100 && !isRequestingNextPage && ListToString(nextPageStoresList).length) {
            if (target?.scrollTop >= (maxScrollLeft - 100)) {
                setIsRequestingNextPage(true)
                handleNextPage()
            }
        }
    }

    useEffect(() => {
        window.history.scrollRestoration = 'manual'
    }, []);

    useEffect(() => {
        if (stores && !products && !isRequestingNextPage) {
            GetProducts({
                params: {
                    "pg": page,
                    "q": search,
                    "pi": minPrice ?? "",
                    "pf": maxPrice ?? "",
                    "stores": ListToString(stores.filter((s) => !storesFilter.includes(s.name))?.map((s) => {
                        return s.name
                    })),
                }
            }).then((r) => {
                if (r) {
                    setProducts(removeDuplicates(r))
                    setIsLoadingProducts(false)
                }
            })
        }
    }, [isRequestingNextPage, maxPrice, minPrice, page, products, search, setIsLoadingProducts, setProducts, stores, storesFilter])

    const cardMinWidth = 350;
    return !isLoadingStores ? !isLoadingProducts ? <SimpleGrid ref={gridRef} id="gridlist" onScroll={handleScroll} overflow={'scroll'} overflowX={'hidden'} justifyItems={"center"} minChildWidth={cardMinWidth} maxH={height - 90} gap={10}>
        {products?.sort((a, b) => parseFloat(a.price) > parseFloat(b.price) ? 1 : -1)?.map((p: Product, i: number) => <ProductCard key={i} code={i} stores={stores} product={p} />)}
        {isRequestingNextPage ?
            <Box mt={100} mb={60}>
                <Text ml={5} fontSize={20} color={"white"}>Indo mais fundo</Text>
                <Spinner ml={15} w={150} h={150} color="#01c0ff" size='xl' />
            </Box>
            :
            <></>}
        <Box h={200}>
        </Box>
    </SimpleGrid>
        :
        <Center>
            <Box>
                <Text fontSize={20} color={"white"}>Carregando Produtos</Text>
                <Spinner ml={15} w={150} h={150} color="#01c0ff" size='xl' />
            </Box>
        </Center>
        :
        <Center>
            <Box>
                <Text fontSize={20} color={"white"}>Carregando Lojas</Text>
                <Spinner ml={15} w={150} h={150} color="#01ff12" size='xl' />
            </Box>
        </Center>
}

export default ProductGrid;