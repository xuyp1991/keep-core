import React from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import * as Icons from "../components/Icons"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useFetchData } from "../hooks/useFetchData"
import { BondingSection } from "../components/BondingSection"

const TBTCApplicationPage = () => {
  // fetch data from service
  const initialData = {}
  const [state] = useFetchData(
    tbtcAuthorizationService.fetchTBTCAuthorizationData,
    initialData
  )

  console.log("state:", state.data)

  const data = [
    {
      operatorAddress: "address",
      stakeAmount: "1000",
      contracts: [
        {
          contractName: "BondedECDSAKeepFactory",
          operatorContractAddress: "address",
          isAuthorized: false,
        },
        {
          contractName: "TBTCSystem",
          operatorContractAddress: "address",
          isAuthorized: true,
        },
      ],
    },
  ]

  // fetch data from service
  const bondingData = [
    {
      operatorAddress: "address",
      stakeAmount: "1000",
      bondedETH: "1000 ETH",
      availableETH: "1000 ETH",
    },
  ]

  return (
    <PageWrapper
      className=""
      title="tBTC"
      nextPageLink="/rewards"
      nextPageTitle="Rewards"
      nextPageIcon={Icons.TBTC}
    >
      <nav className="mb-2">
        <a
          href="https://tbtc.network/"
          className="arrow-link h4"
          rel="noopener noreferrer"
          target="_blank"
        >
          tBTC Website
        </a>
      </nav>
      <AuthorizeContracts data={data} />
      <BondingSection data={bondingData} />
    </PageWrapper>
  )
}

export default TBTCApplicationPage
