import { ChangeEvent, SyntheticEvent, useState, useRef, useEffect } from "react";
import "./SearchForm.css"

import {
    useRive,
    Layout,
    Fit,
    Alignment,
    UseRiveParameters,
    StateMachineInput,
    useStateMachineInput,
    RiveState,
} from "rive-react";

import {
    Box,
    InputGroup,
    Input,
    Center,
    Button,
} from "@chakra-ui/react";
import { useNavigate } from "react-router-dom";
import useWindowDimensions from "../../hooks/useWindowDimensions";
import { IoSearchOutline } from "react-icons/io5";

const STATE_MACHINE_NAME = 'State Machine 1'

const SearchForm = (riveProps: UseRiveParameters = {}) => {
    const [searchValue, setSearchValue] = useState('');
    const [inputLookMultiplier, setInputLookMultiplier] = useState(0);
    const inputRef = useRef(null) as any;
    let navigate = useNavigate();

    useEffect(() => {
        if (inputRef?.current && !inputLookMultiplier) {
            setInputLookMultiplier((inputRef.current.offsetWidth / 100) * 1.25);
        }
    }, [inputLookMultiplier, inputRef])

    const { rive: riveInstance, RiveComponent }: RiveState = useRive({
        src: '/genie.riv',
        stateMachines: STATE_MACHINE_NAME,
        autoplay: true,
        layout: new Layout({
            fit: Fit.Cover,
            alignment: Alignment.Center
        }),
        ...riveProps
    });

    const onSubmit = (e: SyntheticEvent) => {
        navigate("/search/" + searchValue)
        e.preventDefault();
        return false;
    };

    const onEnter = (e: { key: string; }) => {
        if (e.key === 'Enter') {
            navigate("/search/" + searchValue)
        }
    }

    const isCheckingInput: StateMachineInput | any = useStateMachineInput(riveInstance, STATE_MACHINE_NAME, 'isLooking');
    const numLookInput: StateMachineInput | any = useStateMachineInput(riveInstance, STATE_MACHINE_NAME, 'Look');

    useEffect(() => {
        if (isCheckingInput) {
            isCheckingInput!.value = false;
        }
    }, [isCheckingInput])

    const onSearchChange = (e: ChangeEvent<HTMLInputElement>) => {
        const newVal = e.target.value;
        setSearchValue(newVal);

        if (!isCheckingInput!.value) {
            isCheckingInput!.value = true;
        }

        const numChars = newVal.length;
        numLookInput!.value = numChars * inputLookMultiplier;
    };

    const onSearchFocus = () => {
        isCheckingInput!.value = true;
        if (numLookInput!.value !== searchValue.length * inputLookMultiplier) {
            numLookInput!.value = searchValue.length * inputLookMultiplier;
        }
    }

    const onSearchBlur = () => {
        isCheckingInput!.value = false;
    }

    const maxInputSize = (size: number) => {
        if (size < 500) {
            return 300
        } else {
            return 500
        }
    }

    const SearchIcon = IoSearchOutline as React.FC;

    return (
        <Box mt={'9%'}>
            <Box>
                <Center>
                    <RiveComponent className="rive-container" />
                </Center>
                <Center>
                    <Box p={15} maxW={'100%'} backgroundColor={'#161616'} borderColor={'#f1ce5a'} borderStyle={"solid"} borderRadius={15} borderWidth={5}>
                        <InputGroup>
                            <Box w={maxInputSize(useWindowDimensions().width)}>
                                <Center>
                                    <Input
                                        onKeyDown={onEnter}
                                        borderRadius={5}
                                        w={'100%'}
                                        h={35}
                                        fontSize={15}
                                        type="text"
                                        backgroundColor={"gray.300"}
                                        placeholder="Escreva o produto que deseja!"
                                        onFocus={onSearchFocus}
                                        onBlur={onSearchBlur}
                                        value={searchValue}
                                        onChange={onSearchChange}
                                        ref={inputRef}
                                    />
                                </Center>
                                <Center>
                                    <Button
                                        mt={5}
                                        loading={false}
                                        borderRadius={5}
                                        fontWeight={'bold'}
                                        fontSize={15}
                                        color={"white"}
                                        bgColor={"#505050"}
                                        borderWidth={0}
                                        w={"90%"}
                                        h={30}
                                        aria-label="search"
                                        variant='solid'
                                        cursor={'pointer'}
                                        onClick={(event) => {
                                            event.preventDefault();
                                            onSubmit(event);
                                        }}>
                                        <SearchIcon />
                                        Pesquisar
                                    </Button>
                                </Center>
                            </Box>
                        </InputGroup>
                    </Box>
                </Center>
            </Box>
        </Box>

    )
}

export default SearchForm;