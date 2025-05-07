/* eslint-disable @typescript-eslint/no-unused-vars */
import { useParams } from 'react-router-dom';
import { Box } from "@chakra-ui/react"
import GetStores from "../../requests/GetStores";
import Store from "../../types/Store";
import { useEffect, useState, Suspense, lazy, useRef, SetStateAction } from "react";
import ProductHeader from '../../components/ProductHeader';
import FilterModal from '../../components/FilterModal';
import useWindowDimensions from '../../hooks/useWindowDimensions';
import ProductFooter from '../../components/ProductFooter';
import Product from '../../types/Product';
const ProductGrid = lazy(() => import('../../components/ProductGrid'));

const Products = () => {
    const [stores, setStores] = useState<Store[] | undefined>(undefined);
    const [products, setProducts] = useState<Product[] | undefined>(undefined);
    const [isLoadingStores, setIsLoadingStores] = useState(true)
    const [isLoadingProducts, setIsLoadingProducts] = useState(true)
    const [isRequestingNextPage, setIsRequestingNextPage] = useState(false)
    const { search } = useParams();
    const gridRef = useRef(null)

    const [isFilterOpen, setIsFilterOpen] = useState(false)
    const [updateFilters, setUpdateFilters] = useState(false)
    const [forceUpdateFilters, setForceUpdateFilters] = useState(false)
    const [minPrice, setMinPrice] = useState("")
    const [maxPrice, setMaxPrice] = useState("")
    const [storesFilter, setStoresFilter] = useState([""])
    const { width } = useWindowDimensions();
    const [mobileMode, setMobileMode] = useState(false)

    useEffect(() => {
        setIsLoadingStores(true)
        GetStores().then((s: any) => {
            if (s) {
                setStores(s)
                setIsLoadingStores(false)
            }
        })
    }, [])

    useEffect(() => {
        if (width > 750) {
            setMobileMode(true)
        } else {
            setMobileMode(false)
        }
    }, [width])

    return <Box>
        {!isFilterOpen ? <ProductHeader mobileMode={mobileMode} isFilterOpen={isFilterOpen} setIsFilterOpen={setIsFilterOpen} setForceUpdateFilters={setForceUpdateFilters} search={search ?? ""} /> : <></>}
        <Box m={10} pb={200} mt={20}>
            <meta name="referrer" content="no-referrer" />
            <Suspense fallback={<div>Loading...</div>}>
                <ProductGrid
                    isLoadingStores={isLoadingStores}
                    isLoadingProducts={isLoadingProducts}
                    setIsLoadingProducts={setIsLoadingProducts}
                    isRequestingNextPage={isRequestingNextPage}
                    setIsRequestingNextPage={setIsRequestingNextPage}
                    stores={stores}
                    products={products}
                    setProducts={setProducts}
                    search={search ?? ""}
                    minPrice={minPrice}
                    maxPrice={maxPrice}
                    storesFilter={storesFilter}
                    updateFilters={updateFilters}
                    setUpdateFilters={setUpdateFilters}
                    forceUpdateFilters={forceUpdateFilters}
                    setForceUpdateFilters={setForceUpdateFilters}
                    gridRef={gridRef} />
            </Suspense>
        </Box>
        {!mobileMode ? <ProductFooter isFilterOpen={isFilterOpen} setIsFilterOpen={setIsFilterOpen} gridRef={gridRef} /> : <></>}
        <FilterModal
            stores={stores}
            isOpen={isFilterOpen}
            onClose={() => {
                setIsFilterOpen(false)
                setUpdateFilters(true)
            }}
            minPrice={minPrice}
            maxPrice={maxPrice}
            storesFilter={storesFilter}
            setMinPrice={setMinPrice}
            setMaxPrice={setMaxPrice}
            setStoresFilter={setStoresFilter}
            mobileMode={mobileMode}
        />
    </Box>
}

export default Products;