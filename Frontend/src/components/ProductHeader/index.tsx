import { Box, Button, Input } from "@chakra-ui/react";
import useWindowDimensions from "../../hooks/useWindowDimensions";
import { IoSearchOutline } from "react-icons/io5";
import FilterButton from "../FilterButton";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const ProductHeader = ({ mobileMode, isFilterOpen, setIsFilterOpen, search, setForceUpdateFilters }: { mobileMode: boolean, isFilterOpen: boolean, setIsFilterOpen: React.Dispatch<React.SetStateAction<boolean>>, setForceUpdateFilters: React.Dispatch<React.SetStateAction<boolean>>, search: string }) => {
    const [searchString, setSearchString] = useState(search)

    let navigate = useNavigate();

    const handleSearch = () => {
        navigate("/search/" + searchString, { replace: true })
        //window.location.reload()
        setForceUpdateFilters(true)
    }

    const onEnter = (e: { key: string; }) => {
        if (e.key === 'Enter') {
            handleSearch()
        }
    }

    const SearchIcon = IoSearchOutline as React.FC;

    return <Box position={"fixed"} top={0} zIndex={999} display={"flex"} mb={30} pb={2} w={useWindowDimensions().width} h={70} backgroundColor={"#1a181d"}>
        {mobileMode ? <Box w={120} h={"80%"}>
            <FilterButton isFilterOpen={isFilterOpen} setIsFilterOpen={setIsFilterOpen} />
        </Box> : <></>}
        <Box ml={15} w={'100%'} justifyContent={'center'} pt={5} display={"flex"} h={'65%'}>
            <Input backgroundColor={"gray.300"} borderColor={"transparent"} onKeyDown={onEnter} value={searchString} onChange={(e) => setSearchString(e.target.value)} defaultValue={search} maxW={900} fontSize={17} placeholder="Escreva o produto que deseja!" borderLeftRadius={10} w={useWindowDimensions().width * 0.65} size={"sm"} />
            <Button onClick={handleSearch} backgroundColor={"#34a4a5"} borderRightRadius={10} border={"none"} ml={2} fontSize={17} w={50}>
                <SearchIcon />
            </Button>
        </Box>
    </Box>
}

export default ProductHeader;