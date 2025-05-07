import { Box, Text, Button, Center, Image } from "@chakra-ui/react"
import { Modal, ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader, ModalOverlay } from "@chakra-ui/modal"
import { NumberInput, NumberDecrementStepper, NumberIncrementStepper, NumberInputField, NumberInputStepper } from "@chakra-ui/number-input"
import { Tooltip } from "@chakra-ui/tooltip"
import useWindowDimensions from "../../hooks/useWindowDimensions"
import Store from "../../types/Store"
import LazyLoad from "react-lazyload"

import { FaCheck } from "react-icons/fa";
import { RiCloseFill } from "react-icons/ri";

const CheckIcon = FaCheck as React.FC;
const CloseIcon = RiCloseFill as React.FC;

const FilterModal = ({ isOpen, onClose, stores, minPrice, maxPrice, setMinPrice, setMaxPrice, storesFilter, setStoresFilter, mobileMode }:
    {
        isOpen: boolean,
        onClose: () => void,
        stores: Store[] | undefined,
        minPrice: string,
        setMinPrice: React.Dispatch<React.SetStateAction<string>>,
        maxPrice: string,
        setMaxPrice: React.Dispatch<React.SetStateAction<string>>,
        storesFilter: string[],
        setStoresFilter: React.Dispatch<React.SetStateAction<string[]>>,
        mobileMode: boolean
    }) => {

    return <Modal scrollBehavior={"inside"} onClose={onClose} isOpen={isOpen} isCentered>
        <ModalOverlay backdropFilter='blur(10px) hue-rotate(90deg)' />
        <ModalContent margin={"auto"} color={"white"} borderRadius={20} h={"full"} maxW={500} w={'90%'} zIndex={9999}>
            <Box borderRadius={20} mt={!mobileMode ? "5%" : "15%"} mb={"15%"} backgroundColor={"#1a181d"} h={useWindowDimensions().height * 0.8}>
                <ModalHeader borderBottom={'solid'} borderColor={"lightgray"} borderWidth={0.2} p={20} display={'flex'} justifyContent={'space-between'} top={0} h={20}>
                    <Center w={'100%'}><Text fontSize={30}>Filtros</Text> </Center>
                    <ModalCloseButton cursor={"pointer"} borderRadius={5} backgroundColor={"#e8290786"} color={"#ffffff92"} border={"none"} />
                </ModalHeader>
                <ModalBody overflow={'scroll'} overflowX={'hidden'} h={useWindowDimensions().height * 0.57}>
                    <Box mb={30}>
                        <Text m={"auto"} w={"20%"} p={2} fontSize={25}>Preços</Text>
                        <Box display={"flex"} ml={10} pl={4} mb={10} w={'85%'} borderRadius={20} border={"solid"} borderColor={"#0080008d"} borderWidth={2}>
                            <Text m={"auto"} w={"25%"} p={2} mr={20} fontSize={20}>Mínimo:</Text>
                            <Box display={"flex"} w={'85%'}>
                                {minPrice !== "" && minPrice !== undefined ? <Text m={0} mt={15} pr={12} h={10} fontSize={30}>R$</Text> : <></>}
                                <NumberInput
                                    p={2}
                                    display={'flex'}
                                    defaultValue={undefined}
                                    max={10}
                                    keepWithinRange={false}
                                    clampValueOnBlur={false}
                                    value={minPrice}
                                    onChange={(v: any) => {
                                        if (parseFloat(v) > 0 || v === "") {
                                            setMinPrice(v)
                                        } else {
                                            if (parseFloat(v) <= 0) {
                                                setMinPrice("")
                                            }
                                        }
                                    }}
                                >
                                    <NumberInputField color={"white"} outline={"none"} backgroundColor={"transparent"} borderWidth={0} border={"none"} fontSize={30} w={'85%'} h={60} />
                                    <NumberInputStepper color={"white"} pr={10}>
                                        <NumberIncrementStepper />
                                        <NumberDecrementStepper />
                                    </NumberInputStepper>
                                </NumberInput>
                            </Box>
                        </Box>
                        <Box display={"flex"} ml={10} pl={2} w={'85%'} borderRadius={20} border={"solid"} borderColor={"#ff000084"} borderWidth={2}>
                            <Text m={"auto"} w={"25%"} p={2} mr={20} fontSize={20}>Máximo:</Text>
                            <Box display={"flex"} w={'85%'}>
                                {maxPrice !== "" && maxPrice !== undefined ? <Text m={0} mt={15} pr={12} h={10} fontSize={30}>R$</Text> : <></>}
                                <NumberInput
                                    p={2}
                                    display={'flex'}
                                    defaultValue={undefined}
                                    max={10}
                                    keepWithinRange={false}
                                    clampValueOnBlur={false}
                                    value={maxPrice}
                                    onChange={(v: any) => {
                                        if (parseFloat(v) > 0 || v === "") {
                                            setMaxPrice(v)
                                        } else {
                                            if (parseFloat(v) <= 0) {
                                                setMaxPrice("")
                                            }
                                        }
                                    }}
                                >
                                    <NumberInputField color={"white"} outline={"none"} backgroundColor={"transparent"} borderWidth={0} border={"none"} fontSize={30} w={'85%'} h={60} />
                                    <NumberInputStepper color={"white"} pr={10}>
                                        <NumberIncrementStepper />
                                        <NumberDecrementStepper />
                                    </NumberInputStepper>
                                </NumberInput>
                            </Box>
                        </Box>
                    </Box>
                    <Box>
                        <Text pb={10} m={"auto"} w={"20%"} fontSize={25}>Lojas</Text>
                        {stores?.map((store) => {
                            return <Box mb={4} w={"85%"} key={store.name} ml={20} mr={20} display={"flex"} justifyContent={"space-between"} >
                                <LazyLoad height={20} offset={50}>
                                    <Tooltip borderRadius={10} p={10} backgroundColor={"#c6c6c6"} label={store?.name} aria-label='store'>
                                        <Image w={100} objectFit="scale-down" alt={store.name} src={process.env.REACT_APP_SERVER_URL + store.logo} />
                                    </Tooltip>
                                </LazyLoad>
                                {!storesFilter.includes(store.name) ? <Button onClick={() => { setStoresFilter([...storesFilter, store.name]) }} mr={20} h={25} borderColor={"transparent"} color={"green"} borderRadius={5} fontSize={20} backgroundColor={"transparent"}>{<CheckIcon />}</Button>
                                    :
                                    <Button onClick={() => { setStoresFilter(storesFilter.filter((sF) => sF !== store.name)) }} mr={20} h={25} borderColor={"transparent"} color={"red"} borderRadius={5} fontSize={20} backgroundColor={"transparent"}>{<CloseIcon />}</Button>}
                            </Box>
                        })}
                    </Box>
                </ModalBody>
                <ModalFooter justifyContent={'center'} h={'25%'}>
                    <Button cursor={"pointer"} backgroundColor={'#3dd4b3'} fontSize={25} w={'80%'} h={45} border={'none'} borderRadius={20} onClick={onClose}>
                        Aplicar
                    </Button>
                </ModalFooter>
            </Box>
        </ModalContent>
    </Modal>
}

export default FilterModal