import "./ProductCard.css"
import { Box, Button, Center, Container, Image, Link, Spacer } from "@chakra-ui/react";
import { Card, CardBody, CardFooter, CardHeader, } from "@chakra-ui/card";
import { Tooltip } from "@chakra-ui/tooltip"
import LazyLoad from 'react-lazyload';
import { useEffect, useState } from "react";
import {
    useRive,
    Layout,
    Fit,
    Alignment,
    useStateMachineInput,
    RiveState,
} from "rive-react";
import Store from "../../types/Store";
import Product from "../../types/Product";

const STATE_MACHINE_NAME = 'State Machine 1'

const ProductCard = ({ stores, product, code }: { stores: Store[] | undefined, product: Product, code: number }) => {
    const [store, setStore] = useState<Store | undefined>(undefined);

    const currencyFormatter = new Intl.NumberFormat('pt-BR', {
        style: 'currency',
        currency: 'BRL',
    });

    const { rive: riveInstance, RiveComponent }: RiveState = useRive({
        src: '/rating.riv',
        stateMachines: STATE_MACHINE_NAME,
        autoplay: true,
        layout: new Layout({
            fit: Fit.Cover,
            alignment: Alignment.Center
        }),

    });

    const RatingStars = useStateMachineInput(riveInstance, STATE_MACHINE_NAME, 'rating', parseFloat(Number(product.stars).toFixed(1)));

    useEffect(() => {
        if (RatingStars) {
            RatingStars.value = parseFloat(Number(product.stars).toFixed(1));
        }
    }, [RatingStars, product])

    useEffect(() => {
        stores?.map((s) => {
            if (s) {
                if (s.name.toLocaleLowerCase() === product.store.toLocaleLowerCase()) {
                    setStore(s)
                }
            }
            return s
        })
    }, [product.store, stores])



    return <Card className={`card-${code}`} borderRadius={10} w={'100%'} backgroundColor={"#1a181d"}>
        <CardHeader p={5} color={"white"}>
            <Box pt={5} justifyContent={"space-between"} display={"flex"}>
                <Box pt={2} pl={15}>
                    <LazyLoad scroll once scrollContainer={"#gridlist"} height={25} offset={50}>
                        <Tooltip borderRadius={10} p={10} backgroundColor={"#c6c6c6"} label={store?.name} aria-label='store'>
                            {
                                store ? <Image width={140} minH={4} alt={store.name} src={process.env.REACT_APP_SERVER_URL + store.logo} /> : <></>
                            }
                        </Tooltip>
                    </LazyLoad>
                </Box>
                {
                    Number(product.stars).toFixed(1) !== "-1.0" ?
                        <Tooltip borderRadius={10} p={10} backgroundColor={"#c6c6c6"} label={Number(product.stars).toFixed(1)} aria-label='rating'>
                            <Box ml={20} h={30}>
                                <LazyLoad scroll once scrollContainer={"#gridlist"} height={25} offset={50}>
                                    <RiveComponent defaultValue={product.stars} className="rive-rating" />
                                </LazyLoad>
                            </Box>
                        </Tooltip>
                        : <></>
                }
            </Box>
        </CardHeader>
        <CardBody p={10}>
            <Box width={"100%"} h={200} backgroundColor={"white"} cursor={"pointer"}>
                <Center>
                    <LazyLoad scroll scrollContainer={"#gridlist"} once height={20} offset={50}>
                        <Link href={product.link} target="_blank">
                            <Image rel="noopener noreferrer" boxSize={200} objectFit="scale-down" alt={product.name} src={product.image} />
                        </Link>
                    </LazyLoad>
                </Center>
            </Box>
            <Tooltip placement="top" w={350} borderRadius={10} p={10} backgroundColor={"#c6c6c6"} label={product.name} aria-label='name'>
                <Link href={product.link} target="_blank">
                    <Container rel="noopener noreferrer" minHeight={50} fontSize={20} lineHeight={1.5} /*noOfLines={2}*/ pt={5} color={"white"}>
                        {product.name}
                    </Container>
                </Link>
            </Tooltip>
        </CardBody>
        <Spacer />
        <Link w={"100%"} href={product.link} target="_blank">
            <CardFooter w={"100%"} p={5} pb={15} display={"block"} justify={"center"}>
                <Link w={"100%"} href={product.link} target="_blank">
                    <Center w={"100%"} rel="noopener noreferrer" fontSize={27} color={"#06f57e"} p={5}>
                        {currencyFormatter.format(parseFloat(product.price))}
                    </Center>
                </Link>
                <Box w={"100%"} display={"flex"} justifyContent={"center"}>
                    <Button rel="noopener noreferrer" cursor={'pointer'} fontSize={25} border={'none'} backgroundColor={"#34a4a5"} w={"80%"} borderRadius={20}>
                        Acessar
                    </Button>
                </Box>
            </CardFooter>
        </Link>
    </Card>
}

export default ProductCard;