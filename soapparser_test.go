package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseCalendarData(t *testing.T) {
	results := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
   <s:Header>
      <h:ServerVersionInfo xmlns:h="http://schemas.microsoft.com/exchange/services/2006/types" xmlns="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" MajorVersion="14" MinorVersion="3" MajorBuildNumber="210" MinorBuildNumber="2" />
   </s:Header>
   <s:Body xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <m:GetFolderResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
         <m:ResponseMessages>
            <m:GetFolderResponseMessage ResponseClass="Success">
               <m:ResponseCode>NoError</m:ResponseCode>
               <m:Folders>
                  <t:CalendarFolder>
                     <t:FolderId Id="folderid" ChangeKey="changeKey" />
                  </t:CalendarFolder>
               </m:Folders>
            </m:GetFolderResponseMessage>
         </m:ResponseMessages>
      </m:GetFolderResponse>
   </s:Body>
</s:Envelope>`
	item := ParseCalendarFolder(results)

	assert.Equal(t, item.Id, "folderid")
	assert.Equal(t, item.ChangeKey, "changeKey")
}

func TestParseAppointments(t *testing.T) {
	results := `<?xml version="1.0" encoding="UTF-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/">
   <s:Header>
      <h:ServerVersionInfo xmlns:h="http://schemas.microsoft.com/exchange/services/2006/types" xmlns="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" MajorVersion="14" MinorVersion="3" MajorBuildNumber="210" MinorBuildNumber="2" />
   </s:Header>
   <s:Body xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
      <m:FindItemResponse xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types">
         <m:ResponseMessages>
            <m:FindItemResponseMessage ResponseClass="Success">
               <m:ResponseCode>NoError</m:ResponseCode>
               <m:RootFolder TotalItemsInView="2" IncludesLastItemInRange="true">
                  <t:Items>
                     <t:CalendarItem>
                        <t:ItemId Id="firstid" ChangeKey="firstchangekey" />
                        <t:Subject>Travel</t:Subject>
                        <t:DisplayCc />
                        <t:DisplayTo />
                        <t:Start>2015-04-13T16:00:00Z</t:Start>
                        <t:End>2015-04-14T00:00:00Z</t:End>
                        <t:IsAllDayEvent>false</t:IsAllDayEvent>
                        <t:MyResponseType>Organizer</t:MyResponseType>
                        <t:Organizer>
                           <t:Mailbox>
                              <t:Name>Last, First</t:Name>
                           </t:Mailbox>
                        </t:Organizer>
                     </t:CalendarItem>
                     <t:CalendarItem>
                        <t:ItemId Id="secondid" ChangeKey="secondchangekey" />
                        <t:Subject>Coding in Go</t:Subject>
                        <t:DisplayCc>Suthers, Sally; Mr. Smithers</t:DisplayCc>
                        <t:DisplayTo>Legoman; Batman</t:DisplayTo>
                        <t:Start>2015-04-14T04:00:00Z</t:Start>
                        <t:End>2015-04-17T04:00:00Z</t:End>
                        <t:IsAllDayEvent>true</t:IsAllDayEvent>
                        <t:Location>Conference Room 406</t:Location>
                        <t:MyResponseType>Accept</t:MyResponseType>
                        <t:Organizer>
                           <t:Mailbox>
                              <t:Name>Other, Person</t:Name>
                           </t:Mailbox>
                        </t:Organizer>
                     </t:CalendarItem>
                  </t:Items>
               </m:RootFolder>
            </m:FindItemResponseMessage>
         </m:ResponseMessages>
      </m:FindItemResponse>
   </s:Body>
</s:Envelope>`

	appointments := ParseAppointments(results)

	assert.NotNil(t, appointments)
	assert.Equal(t, 2, len(appointments))
	assert.Equal(t, "secondid", appointments[1].ItemId)
	assert.Equal(t, "Coding in Go", appointments[1].Subject)
	assert.Equal(t, "Suthers, Sally; Mr. Smithers", appointments[1].Cc)
	assert.Equal(t, "Legoman; Batman", appointments[1].To)
	assert.Equal(t, true, appointments[1].IsAllDayEvent)
	assert.Equal(t, "Conference Room 406", appointments[1].Location)
	assert.Equal(t, "Accept", appointments[1].MyResponseType)
	assert.Equal(t, "Other, Person", appointments[1].Organizer)

	starttime, _ := time.Parse(time.RFC3339, "2015-04-14T04:00:00Z")
	assert.Equal(t, starttime, appointments[1].Start)
	endtime, _ := time.Parse(time.RFC3339, "2015-04-17T04:00:00Z")
	assert.Equal(t, endtime, appointments[1].End)
}
