import { Button } from "@chakra-ui/react"
import { FaRegEdit } from "react-icons/fa";

const EditIcon = FaRegEdit as React.FC;

const FilterButton = ({ isFilterOpen, setIsFilterOpen }: { isFilterOpen: boolean, setIsFilterOpen: React.Dispatch<React.SetStateAction<boolean>> }) => {

    return <Button onClick={() => setIsFilterOpen(true)} mt={2} ml={15} backgroundColor={"transparent"} borderWidth={2} borderColor={"#34a4a5"} justifyContent={"space-around"} color={"#00ffff"} w={'100%'} h={"100%"} borderRadius={20} fontSize={20}>
        <EditIcon />
        Filtrar
    </Button>
}

export default FilterButton