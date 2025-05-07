import { Box, Button } from "@chakra-ui/react"
import FilterButton from "../FilterButton";
import { FaArrowDown } from "react-icons/fa";

const ArrowUpIcon = FaArrowDown as React.FC;

const ProductFooter = ({ isFilterOpen, setIsFilterOpen, gridRef }:
    { isFilterOpen: boolean, setIsFilterOpen: React.Dispatch<React.SetStateAction<boolean>>, gridRef: React.MutableRefObject<null | HTMLElement> }) => {

    return <Box display={"flex"} justifyContent={"space-around"} position={"fixed"} bottom={0} w={'100%'} backgroundColor={"#1a181d"} h={40}>
        <Box pt={5} w={110} h={20}>
            <FilterButton isFilterOpen={isFilterOpen} setIsFilterOpen={setIsFilterOpen} />
        </Box>
        <Button pt={10} h={20} onClick={() => {
            gridRef.current?.scrollTo({ top: 0, behavior: "smooth" })
        }} cursor={"pointer"} border={"none"} backgroundColor={"transparent"} color={"white"} fontSize={30}>
            <ArrowUpIcon />
        </Button>
    </Box>
}

export default ProductFooter;