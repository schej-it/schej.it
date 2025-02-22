<template>
  <span>
    <div class="tw-select-none tw-py-4" style="-webkit-touch-callout: none">
      <div class="tw-flex tw-flex-col sm:tw-flex-row">
        <div class="tw-flex tw-grow tw-pl-4" :class="isSignUp ? '' : 'tw-pr-4'">
          <template v-if="event.daysOnly">
            <div class="tw-grow">
              <div class="tw-flex tw-items-center tw-justify-between">
                <v-btn
                  :class="hasPrevPage ? 'tw-visible' : 'tw-invisible'"
                  class="tw-border-gray"
                  outlined
                  icon
                  @click="prevPage"
                  ><v-icon>mdi-chevron-left</v-icon></v-btn
                >
                <div
                  class="tw-text-lg tw-font-medium tw-capitalize sm:tw-text-xl"
                >
                  {{ curMonthText }}
                </div>
                <v-btn
                  :class="hasNextPage ? 'tw-visible' : 'tw-invisible'"
                  class="tw-border-gray"
                  outlined
                  icon
                  @click="nextPage"
                  ><v-icon>mdi-chevron-right</v-icon></v-btn
                >
              </div>
              <!-- Header -->
              <div class="tw-flex tw-w-full">
                <div
                  v-for="day in daysOfWeek"
                  class="tw-flex-1 tw-p-2 tw-text-center tw-text-base tw-capitalize tw-text-dark-gray"
                >
                  {{ day }}
                </div>
              </div>
              <!-- Days grid -->
              <div
                id="drag-section"
                class="tw-grid tw-grid-cols-7"
                @mouseleave="resetCurTimeslot"
              >
                <div
                  v-for="(day, i) in monthDays"
                  :key="day.time"
                  class="timeslot tw-aspect-square tw-p-2 tw-text-sm sm:tw-text-base"
                  :class="dayTimeslotClassStyle[i].class"
                  :style="dayTimeslotClassStyle[i].style"
                  v-on="dayTimeslotVon[i]"
                >
                  {{ day.date }}
                </div>
              </div>

              <v-expand-transition>
                <div
                  :key="hintText"
                  v-if="!isPhone && hintTextShown"
                  class="tw-sticky tw-bottom-4 tw-z-10 tw-flex"
                >
                  <div
                    class="tw-mt-2 tw-flex tw-w-full tw-items-center tw-justify-between tw-gap-1 tw-rounded-md tw-bg-off-white tw-p-2 tw-px-[7px] tw-text-sm tw-text-very-dark-gray"
                  >
                    <div class="tw-flex tw-items-center tw-gap-1">
                      <v-icon small>mdi-information-outline</v-icon>
                      {{ hintText }}
                    </div>
                    <v-icon small @click="closeHint">mdi-close</v-icon>
                  </div>
                </div>
              </v-expand-transition>

              <ToolRow
                v-if="!isPhone && !calendarOnly"
                :event="event"
                :state="state"
                :states="states"
                :cur-timezone.sync="curTimezone"
                :show-best-times.sync="showBestTimes"
                :hide-if-needed.sync="hideIfNeeded"
                :is-weekly="isWeekly"
                :calendar-permission-granted="calendarPermissionGranted"
                :week-offset="weekOffset"
                :num-responses="respondents.length"
                :mobile-num-days.sync="mobileNumDays"
                :allow-schedule-event="allowScheduleEvent"
                :show-event-options="showEventOptions"
                :time-type.sync="timeType"
                @toggleShowEventOptions="toggleShowEventOptions"
                @update:weekOffset="(val) => $emit('update:weekOffset', val)"
                @scheduleEvent="scheduleEvent"
                @cancelScheduleEvent="cancelScheduleEvent"
                @confirmScheduleEvent="confirmScheduleEvent"
              />
            </div>
          </template>
          <template v-else>
            <!-- Times -->
            <div
              :class="calendarOnly ? 'tw-w-12' : ''"
              class="tw-w-8 tw-flex-none sm:tw-w-12"
            >
              <div
                :class="calendarOnly ? 'tw-invisible' : 'tw-visible'"
                class="tw-sticky tw-top-14 tw-z-10 -tw-ml-3 tw-mb-3 tw-h-11 tw-bg-white sm:tw-top-16 sm:tw-ml-0"
              >
                <div
                  :class="hasPrevPage ? 'tw-visible' : 'tw-invisible'"
                  class="tw-sticky tw-top-14 tw-ml-0.5 tw-self-start tw-pt-1.5 sm:tw-top-16 sm:-tw-ml-2"
                >
                  <v-btn class="tw-border-gray" outlined icon @click="prevPage"
                    ><v-icon>mdi-chevron-left</v-icon></v-btn
                  >
                </div>
              </div>

              <div
                :class="calendarOnly ? '' : '-tw-ml-3'"
                class="-tw-mt-[8px] sm:tw-ml-0"
              >
                <div
                  v-for="(time, i) in times"
                  :key="i"
                  class="tw-h-4 tw-pr-1 tw-text-right tw-text-xs tw-font-light tw-uppercase sm:tw-pr-2"
                >
                  {{ time.text }}
                </div>
              </div>
            </div>

            <!-- Middle section -->
            <div class="tw-grow">
              <div
                ref="calendar"
                @scroll="onCalendarScroll"
                class="tw-relative tw-flex tw-flex-col"
              >
                <!-- Days -->
                <div
                  :class="
                    sampleCalendarEventsByDay
                      ? undefined
                      : 'tw-sticky tw-top-14'
                  "
                  class="tw-z-10 tw-flex tw-h-14 tw-items-center tw-bg-white sm:tw-top-16"
                >
                  <div
                    v-for="(day, i) in days"
                    :key="i"
                    class="tw-flex-1 tw-bg-white"
                  >
                    <div class="tw-text-center">
                      <div
                        v-if="isSpecificDates || isGroup"
                        class="tw-text-[12px] tw-font-light tw-capitalize tw-text-very-dark-gray sm:tw-text-xs"
                      >
                        {{ day.dateString }}
                      </div>
                      <div class="tw-text-base tw-capitalize sm:tw-text-lg">
                        {{ day.dayText }}
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Calendar -->
                <div class="tw-flex tw-flex-col">
                  <div class="tw-flex-1">
                    <div
                      id="drag-section"
                      data-long-press-delay="500"
                      class="tw-relative tw-flex"
                      @mouseleave="resetCurTimeslot"
                    >
                      <!-- Loader -->
                      <div
                        v-if="showLoader"
                        class="tw-absolute tw-z-10 tw-grid tw-h-full tw-w-full tw-place-content-center"
                      >
                        <v-progress-circular
                          class="tw-text-green"
                          indeterminate
                        />
                      </div>

                      <div
                        v-for="(_, d) in days"
                        :key="d"
                        class="tw-relative tw-flex-1"
                        :class="
                          ((isGroup && loadingCalendarEvents) ||
                            loadingResponses.loading) &&
                          'tw-opacity-50'
                        "
                      >
                        <!-- Timeslots -->
                        <div v-for="(_, t) in times" :key="t" class="tw-w-full">
                          <div
                            class="timeslot tw-h-4"
                            :class="
                              timeslotClassStyle[d * times.length + t]?.class
                            "
                            :style="
                              timeslotClassStyle[d * times.length + t]?.style
                            "
                            v-on="timeslotVon[d * times.length + t]"
                          ></div>
                        </div>

                        <!-- Calendar events -->
                        <div
                          v-if="
                            !loadingCalendarEvents &&
                            (editing ||
                              alwaysShowCalendarEvents ||
                              showCalendarEvents)
                          "
                        >
                          <transition
                            :name="isGroup ? '' : 'fade-transition'"
                            v-for="event in calendarEventsByDay[
                              d + page * maxDaysPerPage
                            ]"
                            :key="event.id"
                            appear
                          >
                            <div
                              class="tw-absolute tw-w-full tw-select-none tw-p-px"
                              :style="{
                                top: `calc(${event.hoursOffset} * 4 * 1rem)`,
                                height: `calc(${event.hoursLength} * 4 * 1rem)`,
                              }"
                              style="pointer-events: none"
                            >
                              <div
                                class="tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-rounded tw-border tw-border-solid tw-p-1 tw-text-xs"
                                :class="
                                  event.free
                                    ? isGroup && !editing
                                      ? 'tw-border-white tw-bg-light-blue tw-opacity-50'
                                      : 'tw-border-dashed tw-border-blue'
                                    : isGroup && !editing
                                    ? 'tw-border-white tw-bg-light-blue'
                                    : 'tw-border-blue'
                                "
                              >
                                <div
                                  :class="`tw-text-${
                                    isGroup &&
                                    state !== states.EDIT_AVAILABILITY
                                      ? 'white'
                                      : noEventNames
                                      ? 'dark-gray'
                                      : 'blue'
                                  }`"
                                  class="ph-no-capture tw-font-medium"
                                >
                                  {{ noEventNames ? "BUSY" : event.summary }}
                                </div>
                              </div>
                            </div>
                          </transition>
                        </div>

                        <!-- Scheduled event -->
                        <div v-if="state === states.SCHEDULE_EVENT">
                          <div
                            v-if="
                              (dragStart && dragStart.col === d) ||
                              (!dragStart &&
                                curScheduledEvent &&
                                curScheduledEvent.dayIndex === d)
                            "
                            class="tw-absolute tw-w-full tw-select-none tw-p-px"
                            :style="scheduledEventStyle"
                            style="pointer-events: none"
                          >
                            <div
                              class="tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-rounded tw-border tw-border-solid tw-border-blue tw-bg-blue tw-p-px tw-text-xs"
                            >
                              <div class="tw-font-medium tw-text-white">
                                {{ event.name }}
                              </div>
                            </div>
                          </div>
                        </div>

                        <!-- Sign up block being dragged -->
                        <div v-if="state === states.EDIT_SIGN_UP_BLOCKS">
                          <div
                            v-if="dragStart && dragStart.col === d"
                            class="tw-absolute tw-w-full tw-select-none tw-p-px"
                            :style="signUpBlockBeingDraggedStyle"
                            style="pointer-events: none"
                          >
                            <SignUpCalendarBlock
                              :title="newSignUpBlockName"
                              titleOnly
                              unsaved
                            />
                          </div>
                        </div>

                        <div v-if="isSignUp">
                          <!-- Sign up blocks -->
                          <div
                            v-for="block in signUpBlocksByDay[
                              d + page * maxDaysPerPage
                            ]"
                            :key="block._id"
                          >
                            <div
                              class="tw-absolute tw-w-full tw-select-none tw-p-px"
                              :style="{
                                top: `calc(${block.hoursOffset} * 4 * 1rem)`,
                                height: `calc(${block.hoursLength} * 4 * 1rem)`,
                              }"
                              @click="handleSignUpBlockClick(block)"
                            >
                              <SignUpCalendarBlock :signUpBlock="block" />
                            </div>
                          </div>

                          <!-- Sign up blocks to be added after hitting 'save' -->
                          <div
                            v-for="block in signUpBlocksToAddByDay[
                              d + page * maxDaysPerPage
                            ]"
                            :key="block._id"
                          >
                            <div
                              class="tw-absolute tw-w-full tw-select-none tw-p-px"
                              :style="{
                                top: `calc(${block.hoursOffset} * 4 * 1rem)`,
                                height: `calc(${block.hoursLength} * 4 * 1rem)`,
                              }"
                            >
                              <SignUpCalendarBlock
                                :title="block.name"
                                titleOnly
                                unsaved
                              />
                            </div>
                          </div>
                        </div>

                        <!-- Overlaid availabilities -->
                        <div v-if="overlayAvailability">
                          <div
                            v-for="(timeBlock, tb) in overlaidAvailability[d]"
                            :key="tb"
                            class="tw-absolute tw-w-full tw-select-none tw-p-px"
                            :style="{
                              top: `calc(${timeBlock.hoursOffset} * 4 * 1rem)`,
                              height: `calc(${timeBlock.hoursLength} * 4 * 1rem)`,
                            }"
                            style="pointer-events: none"
                          >
                            <div
                              class="tw-h-full tw-w-full tw-border-2"
                              :class="
                                timeBlock.type === 'available'
                                  ? 'overlay-avail-shadow-green tw-border-[#00994CB3] tw-bg-[#00994C66]'
                                  : 'overlay-avail-shadow-yellow tw-border-[#997700CC] tw-bg-[#FFE8B8B3]'
                              "
                            ></div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <ZigZag
                  v-if="hasPrevPage"
                  left
                  class="tw-absolute tw-left-0 tw-top-0 tw-h-full tw-w-3"
                />
                <ZigZag
                  v-if="hasNextPage"
                  right
                  class="tw-absolute tw-right-0 tw-top-0 tw-h-full tw-w-3"
                />
              </div>

              <!-- Hint text (desktop) -->
              <v-expand-transition>
                <div
                  :key="hintText"
                  v-if="!isPhone && hintTextShown"
                  class="tw-sticky tw-bottom-4 tw-z-10 tw-flex"
                >
                  <div
                    class="tw-mt-2 tw-flex tw-w-full tw-items-center tw-justify-between tw-gap-1 tw-rounded-md tw-bg-off-white tw-p-2 tw-px-[7px] tw-text-sm tw-text-very-dark-gray"
                  >
                    <div class="tw-flex tw-items-center tw-gap-1">
                      <v-icon small>mdi-information-outline</v-icon>
                      {{ hintText }}
                    </div>
                    <v-icon small @click="closeHint">mdi-close</v-icon>
                  </div>
                </div>
              </v-expand-transition>

              <v-expand-transition>
                <div
                  v-if="
                    state !== states.EDIT_AVAILABILITY &&
                    max !== respondents.length &&
                    Object.keys(fetchedResponses).length !== 0 &&
                    !loadingResponses.loading
                  "
                >
                  <div class="tw-mt-2 tw-text-sm tw-text-dark-gray">
                    Note: There's no time when all
                    {{ respondents.length }} respondents are available.
                  </div>
                </div>
              </v-expand-transition>

              <ToolRow
                v-if="!isPhone && !calendarOnly"
                :event="event"
                :state="state"
                :states="states"
                :cur-timezone.sync="curTimezone"
                :show-best-times.sync="showBestTimes"
                :hide-if-needed.sync="hideIfNeeded"
                :is-weekly="isWeekly"
                :calendar-permission-granted="calendarPermissionGranted"
                :week-offset="weekOffset"
                :num-responses="respondents.length"
                :mobile-num-days.sync="mobileNumDays"
                :allow-schedule-event="allowScheduleEvent"
                :show-event-options="showEventOptions"
                :time-type.sync="timeType"
                @toggleShowEventOptions="toggleShowEventOptions"
                @update:weekOffset="(val) => $emit('update:weekOffset', val)"
                @scheduleEvent="scheduleEvent"
                @cancelScheduleEvent="cancelScheduleEvent"
                @confirmScheduleEvent="confirmScheduleEvent"
              />
            </div>

            <div
              v-if="!calendarOnly"
              :class="calendarOnly ? 'tw-invisible' : 'tw-visible'"
              class="tw-sticky tw-top-14 tw-z-10 tw-mb-4 tw-h-11 tw-bg-white sm:tw-top-16"
            >
              <div
                :class="hasNextPage ? 'tw-visible' : 'tw-invisible'"
                class="tw-sticky tw-top-14 -tw-mr-2 tw-self-start tw-pt-1.5 sm:tw-top-16"
              >
                <v-btn class="tw-border-gray" outlined icon @click="nextPage"
                  ><v-icon>mdi-chevron-right</v-icon></v-btn
                >
              </div>
            </div>
          </template>
        </div>

        <!-- Right hand side content -->

        <div
          v-if="!calendarOnly"
          class="tw-px-4 tw-py-4 sm:tw-sticky sm:tw-top-16 sm:tw-flex-none sm:tw-self-start sm:tw-py-0 sm:tw-pl-0 sm:tw-pr-0 sm:tw-pt-14"
          :style="{ width: rightSideWidth }"
        >
          <!-- Show respondents if not sign up form, otherwise, show sign up blocks -->
          <template v-if="!isSignUp">
            <div
              class="tw-flex tw-flex-col tw-gap-5"
              v-if="state == states.EDIT_AVAILABILITY"
            >
              <div
                v-if="
                  !(
                    calendarPermissionGranted &&
                    !event.daysOnly &&
                    !addingAvailabilityAsGuest
                  )
                "
                class="tw-text-sm tw-italic tw-text-dark-gray"
              >
                {{
                  (userHasResponded && !addingAvailabilityAsGuest) || curGuestId
                    ? "Editing"
                    : "Adding"
                }}
                availability as
                {{
                  authUser && !addingAvailabilityAsGuest
                    ? `${authUser.firstName} ${authUser.lastName}`
                    : curGuestId?.length > 0
                    ? curGuestId
                    : "a guest"
                }}
              </div>
              <AvailabilityTypeToggle
                v-if="!isGroup && !isPhone"
                class="tw-w-full"
                v-model="availabilityType"
              />
              <!-- User's calendar accounts -->
              <CalendarAccounts
                v-if="
                  calendarPermissionGranted &&
                  !event.daysOnly &&
                  !addingAvailabilityAsGuest
                "
                :toggleState="true"
                :eventId="event._id"
                :calendar-events-map="calendarEventsMap"
                :syncWithBackend="!isGroup"
                :allowAddCalendarAccount="!isGroup"
                @toggleCalendarAccount="toggleCalendarAccount"
                @toggleSubCalendarAccount="toggleSubCalendarAccount"
                :initialCalendarAccountsData="
                  isGroup ? sharedCalendarAccounts : authUser.calendarAccounts
                "
              ></CalendarAccounts>

              <div v-if="showOverlayAvailabilityToggle">
                <v-switch
                  id="overlay-availabilities-toggle"
                  inset
                  :input-value="overlayAvailability"
                  @change="updateOverlayAvailability"
                  hide-details
                >
                  <template v-slot:label>
                    <div class="tw-text-sm tw-text-black">
                      Overlay availabilities
                    </div>
                  </template>
                </v-switch>

                <div class="tw-mt-2 tw-text-xs tw-text-dark-gray">
                  View everyone's availability while inputting your own
                </div>
              </div>

              <!-- Options section -->
              <div
                v-if="!event.daysOnly && showCalendarOptions"
                ref="optionsSection"
              >
                <ExpandableSection
                  label="Options"
                  :value="showEditOptions"
                  @input="toggleShowEditOptions"
                >
                  <div class="tw-flex tw-flex-col tw-gap-5 tw-pt-2.5">
                    <v-dialog
                      v-if="showCalendarOptions"
                      v-model="calendarOptionsDialog"
                      width="500"
                    >
                      <template v-slot:activator="{ on, attrs }">
                        <v-btn
                          outlined
                          class="tw-border-gray tw-text-sm"
                          v-on="on"
                          v-bind="attrs"
                        >
                          Calendar options...
                        </v-btn>
                      </template>

                      <v-card>
                        <v-card-title class="tw-flex">
                          <div>Calendar options</div>
                          <v-spacer />
                          <v-btn icon @click="calendarOptionsDialog = false">
                            <v-icon>mdi-close</v-icon>
                          </v-btn>
                        </v-card-title>
                        <v-card-text
                          class="tw-flex tw-flex-col tw-gap-6 tw-pb-8 tw-pt-2"
                        >
                          <AlertText v-if="isGroup" class="-tw-mb-4">
                            Calendar options will only updated for the current
                            group
                          </AlertText>

                          <BufferTimeSwitch
                            :bufferTime.sync="bufferTime"
                            :syncWithBackend="!isGroup"
                          />

                          <WorkingHoursToggle
                            :workingHours.sync="workingHours"
                            :timezone="curTimezone"
                            :syncWithBackend="!isGroup"
                          />
                        </v-card-text>
                      </v-card>
                    </v-dialog>
                  </div>
                </ExpandableSection>
              </div>

              <!-- Delete availability button -->
              <div
                v-if="
                  (!addingAvailabilityAsGuest && userHasResponded) || curGuestId
                "
              >
                <v-dialog
                  v-model="deleteAvailabilityDialog"
                  width="500"
                  persistent
                >
                  <template v-slot:activator="{ on, attrs }">
                    <span
                      v-bind="attrs"
                      v-on="on"
                      class="tw-cursor-pointer tw-text-sm tw-text-red"
                    >
                      {{ !isGroup ? "Delete availability" : "Leave group" }}
                    </span>
                  </template>

                  <v-card>
                    <v-card-title>Are you sure?</v-card-title>
                    <v-card-text class="tw-text-sm tw-text-dark-gray"
                      >Are you sure you want to
                      {{
                        !isGroup
                          ? "delete your availability from this event?"
                          : "leave this group?"
                      }}</v-card-text
                    >
                    <v-card-actions>
                      <v-spacer />
                      <v-btn text @click="deleteAvailabilityDialog = false"
                        >Cancel</v-btn
                      >
                      <v-btn
                        text
                        color="error"
                        @click="
                          $emit('deleteAvailability')
                          deleteAvailabilityDialog = false
                        "
                        >{{ !isGroup ? "Delete" : "Leave" }}</v-btn
                      >
                    </v-card-actions>
                  </v-card>
                </v-dialog>
              </div>
            </div>
            <template v-else>
              <RespondentsList
                ref="respondentsList"
                :event="event"
                :eventId="event._id"
                :days="allDays"
                :times="times"
                :curDate="getDateFromRowCol(curTimeslot.row, curTimeslot.col)"
                :curRespondent="curRespondent"
                :curRespondents="curRespondents"
                :curTimeslot="curTimeslot"
                :curTimeslotAvailability="curTimeslotAvailability"
                :respondents="respondents"
                :parsedResponses="parsedResponses"
                :isOwner="isOwner"
                :isGroup="isGroup"
                :attendees="event.attendees"
                :showCalendarEvents.sync="showCalendarEvents"
                :responsesFormatted="responsesFormatted"
                :timezone="curTimezone"
                :show-best-times.sync="showBestTimes"
                :hide-if-needed.sync="hideIfNeeded"
                :start-calendar-on-monday.sync="startCalendarOnMonday"
                :show-event-options="showEventOptions"
                :guestAddedAvailability="guestAddedAvailability"
                :addingAvailabilityAsGuest="addingAvailabilityAsGuest"
                @toggleShowEventOptions="toggleShowEventOptions"
                @addAvailability="$emit('addAvailability')"
                @addAvailabilityAsGuest="$emit('addAvailabilityAsGuest')"
                @mouseOverRespondent="mouseOverRespondent"
                @mouseLeaveRespondent="mouseLeaveRespondent"
                @clickRespondent="clickRespondent"
                @editGuestAvailability="editGuestAvailability"
                @refreshEvent="refreshEvent"
              />
            </template>
          </template>
          <template v-else>
            <div class="tw-mb-2 tw-text-lg tw-text-black">Slots</div>
            <div v-if="!isOwner" class="tw-mb-3 tw-flex tw-flex-col">
              <div
                class="tw-flex tw-flex-col tw-gap-1 tw-rounded-md tw-bg-light-gray tw-p-3 tw-text-xs tw-italic tw-text-dark-gray"
              >
                <div v-if="!authUser || alreadyRespondedToSignUpForm">
                  <a class="tw-underline" :href="`mailto:${event.ownerId}`"
                    >Contact sign up creator</a
                  >
                  to edit your slot
                </div>
                <div v-if="event.blindAvailabilityEnabled">
                  Responses are only visible to creator
                </div>
              </div>
            </div>
            <SignUpBlocksList
              ref="signUpBlocksList"
              :signUpBlocks="signUpBlocksByDay.flat()"
              :signUpBlocksToAdd="signUpBlocksToAddByDay.flat()"
              :isEditing="state == states.EDIT_SIGN_UP_BLOCKS"
              :isOwner="isOwner"
              :alreadyResponded="alreadyRespondedToSignUpForm"
              :anonymous="event.blindAvailabilityEnabled"
              @update:signUpBlock="editSignUpBlock"
              @delete:signUpBlock="deleteSignUpBlock"
              @signUpForBlock="$emit('signUpForBlock', $event)"
            />
          </template>
        </div>
      </div>

      <ToolRow
        v-if="isPhone && !calendarOnly"
        class="tw-px-4"
        :event="event"
        :state="state"
        :states="states"
        :cur-timezone.sync="curTimezone"
        :show-best-times.sync="showBestTimes"
        :hide-if-needed.sync="hideIfNeeded"
        :start-calendar-on-monday.sync="startCalendarOnMonday"
        :is-weekly="isWeekly"
        :calendar-permission-granted="calendarPermissionGranted"
        :week-offset="weekOffset"
        :num-responses="respondents.length"
        :mobile-num-days.sync="mobileNumDays"
        :allow-schedule-event="allowScheduleEvent"
        :show-event-options="showEventOptions"
        :time-type.sync="timeType"
        @toggleShowEventOptions="toggleShowEventOptions"
        @update:weekOffset="(val) => $emit('update:weekOffset', val)"
        @scheduleEvent="scheduleEvent"
        @cancelScheduleEvent="cancelScheduleEvent"
        @confirmScheduleEvent="confirmScheduleEvent"
      />

      <!-- Fixed bottom section for mobile -->
      <div
        v-if="isPhone && !calendarOnly"
        class="tw-fixed tw-bottom-16 tw-z-20 tw-w-full"
      >
        <!-- Hint text (mobile) -->
        <v-expand-transition>
          <template v-if="hintTextShown">
            <div :key="hintText">
              <div
                :class="`tw-flex tw-w-full tw-items-center tw-justify-between tw-gap-1 tw-bg-light-gray tw-px-2 tw-py-2 tw-text-sm tw-text-very-dark-gray`"
              >
                <div :class="`tw-flex tw-gap-${hintText.length > 60 ? 2 : 1}`">
                  <v-icon small>mdi-information-outline</v-icon>
                  <div>
                    {{ hintText }}
                  </div>
                </div>
                <v-icon small @click="closeHint">mdi-close</v-icon>
              </div>
            </div>
          </template>
        </v-expand-transition>

        <!-- Fixed pos availability toggle (mobile) -->
        <v-expand-transition>
          <div v-if="!isGroup && editing && !isSignUp">
            <div class="tw-bg-white tw-p-4">
              <AvailabilityTypeToggle
                class="tw-w-full"
                v-model="availabilityType"
              />
            </div>
          </div>
        </v-expand-transition>

        <!-- GCal week selector -->
        <v-expand-transition>
          <div v-if="isWeekly && editing && calendarPermissionGranted">
            <div class="tw-h-16 tw-text-sm">
              <GCalWeekSelector
                :week-offset="weekOffset"
                @update:weekOffset="(val) => $emit('update:weekOffset', val)"
                :start-on-monday="event.startOnMonday"
              />
            </div>
          </div>
        </v-expand-transition>

        <!-- Respondents list -->
        <v-expand-transition>
          <div v-if="delayedShowStickyRespondents">
            <div class="tw-bg-white tw-p-4">
              <RespondentsList
                :max-height="100"
                :event="event"
                :eventId="event._id"
                :days="allDays"
                :times="times"
                :curDate="getDateFromRowCol(curTimeslot.row, curTimeslot.col)"
                :curRespondent="curRespondent"
                :curRespondents="curRespondents"
                :curTimeslot="curTimeslot"
                :curTimeslotAvailability="curTimeslotAvailability"
                :respondents="respondents"
                :parsedResponses="parsedResponses"
                :isOwner="isOwner"
                :isGroup="isGroup"
                :attendees="event.attendees"
                :showCalendarEvents.sync="showCalendarEvents"
                :responsesFormatted="responsesFormatted"
                :timezone="curTimezone"
                :show-best-times.sync="showBestTimes"
                :hide-if-needed.sync="hideIfNeeded"
                :show-event-options="showEventOptions"
                :guestAddedAvailability="guestAddedAvailability"
                :addingAvailabilityAsGuest="addingAvailabilityAsGuest"
                @toggleShowEventOptions="toggleShowEventOptions"
                @addAvailability="$emit('addAvailability')"
                @addAvailabilityAsGuest="$emit('addAvailabilityAsGuest')"
                @mouseOverRespondent="mouseOverRespondent"
                @mouseLeaveRespondent="mouseLeaveRespondent"
                @clickRespondent="clickRespondent"
                @editGuestAvailability="editGuestAvailability"
                @refreshEvent="refreshEvent"
              />
            </div>
          </div>
        </v-expand-transition>
      </div>
    </div>
  </span>
</template>

<style scoped>
.animate-bg-color {
  transition: background-color 0.25s ease-in-out;
}

.break {
  flex-basis: 100%;
  height: 0;
}
</style>

<style>
/* Make timezone select element the same width as content */
#timezone-select {
  width: 5px;
}
</style>

<script>
import {
  timeNumToTimeText,
  dateCompare,
  getDateHoursOffset,
  post,
  put,
  isBetween,
  clamp,
  isPhone,
  utcTimeToLocalTime,
  splitTimeBlocksByDay,
  getTimeBlock,
  dateToDowDate,
  _delete,
  get,
  getDateDayOffset,
  isDateBetween,
  generateEnabledCalendarsPayload,
  isTouchEnabled,
  isElementInViewport,
  lightOrDark,
  removeTransparencyFromHex,
  userPrefers12h,
  getCalendarAccountKey,
  getISODateString,
  getDateWithTimezone,
  timeNumToTimeString,
} from "@/utils"
import {
  availabilityTypes,
  calendarOptionsDefaults,
  eventTypes,
  timeTypes,
} from "@/constants"
import { mapMutations, mapActions, mapState } from "vuex"
import UserAvatarContent from "@/components/UserAvatarContent.vue"
import CalendarAccounts from "@/components/settings/CalendarAccounts.vue"
import Advertisement from "@/components/event/Advertisement.vue"
import SignUpBlock from "@/components/sign_up_form/SignUpBlock.vue"
import SignUpCalendarBlock from "@/components/sign_up_form/SignUpCalendarBlock.vue"
import SignUpBlocksList from "@/components/sign_up_form/SignUpBlocksList.vue"
import ZigZag from "./ZigZag.vue"
import ConfirmDetailsDialog from "./ConfirmDetailsDialog.vue"
import ToolRow from "./ToolRow.vue"
import RespondentsList from "./RespondentsList.vue"
import GCalWeekSelector from "./GCalWeekSelector.vue"
import ExpandableSection from "../ExpandableSection.vue"
import WorkingHoursToggle from "./WorkingHoursToggle.vue"
import AlertText from "../AlertText.vue"

import dayjs from "dayjs"
import ObjectID from "bson-objectid"
import utcPlugin from "dayjs/plugin/utc"
import timezonePlugin from "dayjs/plugin/timezone"
import AvailabilityTypeToggle from "./AvailabilityTypeToggle.vue"
import BufferTimeSwitch from "./BufferTimeSwitch.vue"
dayjs.extend(utcPlugin)
dayjs.extend(timezonePlugin)

export default {
  name: "ScheduleOverlap",
  props: {
    event: { type: Object, required: true },

    loadingCalendarEvents: { type: Boolean, default: false }, // Whether we are currently loading the calendar events
    calendarEventsMap: { type: Object, default: () => {} }, // Object of different users' calendar events
    sampleCalendarEventsByDay: { type: Array, required: false }, // Sample calendar events to use for example calendars
    calendarPermissionGranted: { type: Boolean, default: false }, // Whether user has granted google calendar permissions

    weekOffset: { type: Number, default: 0 }, // Week offset used for displaying calendar events on weekly schejs

    alwaysShowCalendarEvents: { type: Boolean, default: false }, // Whether to show calendar events all the time
    noEventNames: { type: Boolean, default: false }, // Whether to show "busy" instead of the event name
    calendarOnly: { type: Boolean, default: false }, // Whether to only show calendar and not respondents or any other controls
    interactable: { type: Boolean, default: true }, // Whether to allow user to interact with component
    showSnackbar: { type: Boolean, default: true }, // Whether to show snackbar when availability is automatically filled in
    animateTimeslotAlways: { type: Boolean, default: false }, // Whether to animate timeslots all the time
    showHintText: { type: Boolean, default: true }, // Whether to show the hint text telling user what to do

    curGuestId: { type: String, default: "" }, // Id of the current guest being edited
    addingAvailabilityAsGuest: { type: Boolean, default: false }, // Whether the signed in user is adding availability as a guest

    initialTimezone: { type: Object, default: () => ({}) },

    // Availability Groups
    calendarAvailabilities: { type: Object, default: () => ({}) },
  },
  data() {
    return {
      states: {
        HEATMAP: "heatmap", // Display heatmap of availabilities
        SINGLE_AVAILABILITY: "single_availability", // Show one person's availability
        SUBSET_AVAILABILITY: "subset_availability", // Show availability for a subset of people
        BEST_TIMES: "best_times", // Show only the times that work for most people
        EDIT_AVAILABILITY: "edit_availability", // Edit current user's availability
        EDIT_SIGN_UP_BLOCKS: "edit_sign_up_blocks", // Edit the slots on a sign up form
        SCHEDULE_EVENT: "schedule_event", // Schedule event on gcal
      },
      state: "best_times",

      availability: new Set(), // The current user's availability
      ifNeeded: new Set(), // The current user's "if needed" availability
      availabilityAnimTimeouts: [], // Timeouts for availability animation
      availabilityAnimEnabled: false, // Whether to animate timeslots changing colors
      maxAnimTime: 1200, // Max amount of time for availability animation
      unsavedChanges: false, // If there are unsaved availability changes
      curTimeslot: { row: -1, col: -1 }, // The currently highlighted timeslot
      timeslotSelected: false, // Whether a timeslot is selected (used to persist selection on desktop)
      curTimeslotAvailability: {}, // The users available for the current timeslot
      curRespondent: "", // Id of the active respondent (set on hover)
      curRespondents: [], // Id of currently selected respondents (set on click)
      sharedCalendarAccounts: {}, // The user's calendar accounts for changing calendar options for groups
      fetchedResponses: {}, // Responses fetched from the server for the dates currently shown
      loadingResponses: { loading: false, lastFetched: new Date().getTime() }, // Whether we're currently fetching the responses
      responsesFormatted: new Map(), // Map where date/time is mapped to the people that are available then

      /* Sign up form */
      signUpBlocksByDay: [], // The current event's sign up blocks by day
      signUpBlocksToAddByDay: [], // The sign up blocks to be added after hitting 'save'

      /* Edit options */
      showEditOptions:
        localStorage["showEditOptions"] == undefined
          ? false
          : localStorage["showEditOptions"] == "true",
      availabilityType: availabilityTypes.AVAILABLE, // The current availability type
      overlayAvailability: false, // Whether to overlay everyone's availability when editing
      bufferTime: calendarOptionsDefaults.bufferTime, // Set in mounted()
      workingHours: calendarOptionsDefaults.workingHours, // Set in mounted()

      /* Event Options */
      showEventOptions:
        localStorage["showEventOptions"] == undefined
          ? false
          : localStorage["showEventOptions"] == "true",
      showBestTimes:
        localStorage["showBestTimes"] == undefined
          ? false
          : localStorage["showBestTimes"] == "true",
      hideIfNeeded: false,

      /* Variables for drag stuff */
      DRAG_TYPES: {
        ADD: "add",
        REMOVE: "remove",
      },
      timeslot: {
        width: 0,
        height: 0,
      },
      dragging: false,
      dragType: "add",
      dragStart: null,
      dragCur: null,

      /* Variables for options */
      curTimezone: this.initialTimezone,
      curScheduledEvent: null, // The scheduled event represented in the form {hoursOffset, hoursLength, dayIndex}
      timeType:
        localStorage["timeType"] ??
        (userPrefers12h() ? timeTypes.HOUR12 : timeTypes.HOUR24), // Whether 12-hour or 24-hour
      showCalendarEvents: false,
      startCalendarOnMonday: false,
      // localStorage["startCalendarOnMonday"] == undefined
      //   ? false
      //   : localStorage["startCalendarOnMonday"] == "true",

      /* Dialogs */
      deleteAvailabilityDialog: false,
      calendarOptionsDialog: false,

      /* Variables for scrolling */
      optionsVisible: false,
      calendarScrollLeft: 0, // The current scroll position of the calendar
      calendarMaxScroll: 0, // The maximum scroll amount of the calendar, scrolling to this point means we have scrolled to the end
      scrolledToRespondents: false, // whether we have scrolled to the respondents section
      delayedShowStickyRespondents: false, // showStickyRespondents variable but changes 100ms after the actual variable changes (to add some delay)
      delayedShowStickyRespondentsTimeout: null, // Timeout that sets delayedShowStickyRespondents

      /* Variables for pagination */
      page: 0,
      mobileNumDays: localStorage["mobileNumDays"]
        ? parseInt(localStorage["mobileNumDays"])
        : 3, // The number of days to show at a time on mobile
      pageHasChanged: false,

      hasRefreshedAuthUser: false,

      /* Variables for hint */
      hintState: true,

      /** Groups */
      manualAvailability: {},

      /** Constants */
      months: [
        "jan",
        "feb",
        "mar",
        "apr",
        "may",
        "jun",
        "jul",
        "aug",
        "sep",
        "oct",
        "nov",
        "dec",
      ],
    }
  },
  computed: {
    ...mapState(["authUser", "overlayAvailabilitiesEnabled"]),
    /** Returns the width of the right side of the calendar */
    rightSideWidth() {
      if (this.isPhone) return "100%"
      return this.isSignUp ? "18rem" : "13rem"
    },
    /** Returns the days of the week in the correct order */
    daysOfWeek() {
      return !this.startCalendarOnMonday
        ? ["sun", "mon", "tue", "wed", "thu", "fri", "sat"]
        : ["mon", "tue", "wed", "thu", "fri", "sat", "sun"]
    },
    /** Only allow scheduling when a curScheduledEvent exists */
    allowScheduleEvent() {
      return !!this.curScheduledEvent
    },
    /** Returns the availability as an array */
    availabilityArray() {
      return [...this.availability].map((item) => new Date(item))
    },
    /** Returns the if needed availability as an array */
    ifNeededArray() {
      return [...this.ifNeeded].map((item) => new Date(item))
    },
    allowDrag() {
      return (
        this.state === this.states.EDIT_AVAILABILITY ||
        this.state === this.states.EDIT_SIGN_UP_BLOCKS ||
        this.state === this.states.SCHEDULE_EVENT
      )
    },
    /** Returns an array of calendar events for all of the authUser's enabled calendars, separated by the day they occur on */
    calendarEventsByDay() {
      // If this is an example calendar
      if (this.sampleCalendarEventsByDay) return this.sampleCalendarEventsByDay

      // If the user isn't logged in or is adding availability as a guest
      if (!this.authUser || this.addingAvailabilityAsGuest) return []

      let events = []
      let event

      const calendarAccounts = this.isGroup
        ? this.sharedCalendarAccounts
        : this.authUser.calendarAccounts

      // Adds events from calendar accounts that are enabled
      for (const id in calendarAccounts) {
        if (!calendarAccounts[id].enabled) continue

        if (this.calendarEventsMap.hasOwnProperty(id)) {
          for (const index in this.calendarEventsMap[id].calendarEvents) {
            event = this.calendarEventsMap[id].calendarEvents[index]

            // Check if we need to update authUser (to get latest subcalendars)
            const subCalendars = calendarAccounts[id].subCalendars
            if (!subCalendars || !(event.calendarId in subCalendars)) {
              // authUser doesn't contain the subCalendar, so push event to events without checking if subcalendar is enabled
              // and queue the authUser to be refreshed
              events.push(event)
              if (!this.hasRefreshedAuthUser && !this.isGroup) {
                this.refreshAuthUser()
              }
              continue
            }

            // Push event to events if subcalendar is enabled
            if (subCalendars[event.calendarId].enabled) {
              events.push(event)
            }
          }
        }
      }

      const eventsCopy = JSON.parse(JSON.stringify(events))

      const calendarEventsByDay = splitTimeBlocksByDay(
        this.event,
        eventsCopy,
        this.weekOffset
      )

      return calendarEventsByDay
    },
    /** [SPECIFIC TO GROUPS] Returns an object mapping user ids to their calendar events separated by the day they occur on */
    groupCalendarEventsByDay() {
      if (this.event.type !== eventTypes.GROUP) return {}

      const userIdToEventsByDay = {}
      for (const userId in this.event.responses) {
        if (userId === this.authUser._id) {
          userIdToEventsByDay[userId] = this.calendarEventsByDay
        } else if (userId in this.calendarAvailabilities) {
          userIdToEventsByDay[userId] = splitTimeBlocksByDay(
            this.event,
            this.calendarAvailabilities[userId],
            this.weekOffset
          )
        }
      }

      return userIdToEventsByDay
    },
    curRespondentsSet() {
      return new Set(this.curRespondents)
    },

    // -----------------------------------
    //#region Sign up form
    // -----------------------------------

    /** Returns the name of the new sign up block being dragged */
    newSignUpBlockName() {
      return `Slot #${
        this.signUpBlocksByDay.flat().length +
        this.signUpBlocksToAddByDay.flat().length +
        1
      }`
    },

    /** Returns the max allowable drag */
    maxSignUpBlockRowSize() {
      if (!this.dragStart) return null

      const selectedDay = this.signUpBlocksByDay[this.dragStart.col]
      const selectedDayToAdd = this.signUpBlocksToAddByDay[this.dragStart.col]

      if (selectedDay.length === 0 && selectedDayToAdd.length === 0) return null

      let maxSize = Infinity
      for (const block of [...selectedDay, ...selectedDayToAdd]) {
        if (block.hoursOffset * 4 > this.dragStart.row) {
          maxSize = Math.min(
            maxSize,
            block.hoursOffset * 4 - this.dragStart.row
          )
        }
      }

      return maxSize
    },

    /** Whether the current user has already responded to the sign up form */
    alreadyRespondedToSignUpForm() {
      if (!this.authUser || !this.signUpBlocksByDay) return false

      return this.signUpBlocksByDay.some((dayBlocks) =>
        dayBlocks.some((block) =>
          block.responses?.some(
            (response) => response.userId === this.authUser._id
          )
        )
      )
    },

    //#endregion

    /** Returns the max number of people in the curRespondents array available at any given time */
    curRespondentsMax() {
      let max = 0
      if (this.event.daysOnly) {
        for (const day of this.allDays) {
          const num = [
            ...(this.responsesFormatted.get(day.dateObject.getTime()) ??
              new Set()),
          ].filter((r) => this.curRespondentsSet.has(r)).length

          if (num > max) max = num
        }
      } else {
        for (const day of this.allDays) {
          for (const time of this.times) {
            const num = [
              ...this.getRespondentsForHoursOffset(
                day.dateObject,
                time.hoursOffset
              ),
            ].filter((r) => this.curRespondentsSet.has(r)).length

            if (num > max) max = num
          }
        }
      }
      return max
    },
    /** Returns the day offset caused by the timezone offset. If the timezone offset changes the date, dayOffset != 0 */
    dayOffset() {
      return Math.floor((this.event.startTime - this.timezoneOffset / 60) / 24)
    },
    /** Returns all the days that are encompassed by startDate and endDate */
    allDays() {
      const days = []

      for (let i = 0; i < this.event.dates.length; ++i) {
        const date = new Date(this.event.dates[i])
        const offsetDate = new Date(date)
        offsetDate.setDate(offsetDate.getDate() + this.dayOffset)

        let dateString = ""
        if (this.isSpecificDates) {
          dateString = `${
            this.months[offsetDate.getUTCMonth()]
          } ${offsetDate.getUTCDate()}`
        } else if (this.isGroup) {
          const tmpDate = dateToDowDate(
            this.event.dates,
            offsetDate,
            this.weekOffset,
            true
          )

          dateString = `${
            this.months[tmpDate.getUTCMonth()]
          } ${tmpDate.getUTCDate()}`
        }

        days.push({
          dayText: this.daysOfWeek[offsetDate.getUTCDay()],
          dateString,
          dateObject: date,
        })
      }

      return days
    },
    /** Returns a subset of all days based on the page number */
    days() {
      return this.allDays.slice(
        this.page * this.maxDaysPerPage,
        Math.min(this.event.dates.length, (this.page + 1) * this.maxDaysPerPage)
      )
    },
    /** Returns all the days of the month */
    monthDays() {
      const monthDays = []
      const allDaysSet = new Set(
        this.allDays.map((d) => d.dateObject.getTime())
      )

      // Calculate monthIndex and year from event start date and page num
      const date = new Date(this.event.dates[0])
      const monthIndex = date.getUTCMonth() + this.page
      const year = date.getUTCFullYear()

      const lastDayOfPrevMonth = new Date(Date.UTC(year, monthIndex, 0))
      const lastDayOfCurMonth = new Date(Date.UTC(year, monthIndex + 1, 0))

      // Calculate num days from prev month, cur month, and next month to show
      const curDate = new Date(lastDayOfPrevMonth)
      let numDaysFromPrevMonth = 0
      const numDaysInCurMonth = lastDayOfCurMonth.getUTCDate()
      const numDaysFromNextMonth = 6 - lastDayOfCurMonth.getUTCDay()
      const hasDaysFromPrevMonth = !this.startCalendarOnMonday
        ? lastDayOfPrevMonth.getUTCDay() < 6
        : lastDayOfPrevMonth.getUTCDay() != 0
      if (hasDaysFromPrevMonth) {
        curDate.setUTCDate(
          curDate.getUTCDate() -
            (lastDayOfPrevMonth.getUTCDay() -
              (this.startCalendarOnMonday ? 1 : 0))
        )
        numDaysFromPrevMonth = lastDayOfPrevMonth.getUTCDay() + 1
      } else {
        curDate.setUTCDate(curDate.getUTCDate() + 1)
      }
      curDate.setUTCHours(this.event.startTime)

      // Add all days from prev month, cur month, and next month
      const totalDays =
        numDaysFromPrevMonth + numDaysInCurMonth + numDaysFromNextMonth
      for (let i = 0; i < totalDays; ++i) {
        // Only include days from the current month
        if (curDate.getUTCMonth() === lastDayOfCurMonth.getUTCMonth()) {
          monthDays.push({
            date: curDate.getUTCDate(),
            time: curDate.getTime(),
            dateObject: new Date(curDate),
            included: allDaysSet.has(curDate.getTime()),
          })
        } else {
          monthDays.push({
            date: "",
            time: curDate.getTime(),
            dateObject: new Date(curDate),
            included: false,
          })
        }

        curDate.setUTCDate(curDate.getUTCDate() + 1)
      }

      return monthDays
    },
    /** Map from datetime to whether that month day is included  */
    monthDayIncluded() {
      const includedMap = new Map()
      for (const monthDay of this.monthDays) {
        includedMap.set(monthDay.dateObject.getTime(), monthDay.included)
      }
      return includedMap
    },
    /** Returns the text to show for the current month */
    curMonthText() {
      const date = new Date(this.event.dates[0])
      const monthIndex = date.getUTCMonth() + this.page
      const year = date.getUTCFullYear()
      const lastDayOfCurMonth = new Date(Date.UTC(year, monthIndex + 1, 0))

      const monthText = this.months[lastDayOfCurMonth.getUTCMonth()]
      const yearText = lastDayOfCurMonth.getUTCFullYear()
      return `${monthText} ${yearText}`
    },
    defaultState() {
      // Either the heatmap or the best_times state, depending on the toggle
      return this.showBestTimes ? this.states.BEST_TIMES : this.states.HEATMAP
    },
    editing() {
      // Returns whether currently in the editing state
      return (
        this.state === this.states.EDIT_AVAILABILITY ||
        this.state === this.states.EDIT_SIGN_UP_BLOCKS
      )
    },
    scheduling() {
      // Returns whether currently in the scheduling state
      return this.state === this.states.SCHEDULE_EVENT
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
    isOwner() {
      return this.authUser?._id === this.event.ownerId
    },
    isSpecificDates() {
      return this.event.type === eventTypes.SPECIFIC_DATES || !this.event.type
    },
    isWeekly() {
      return this.event.type === eventTypes.DOW
    },
    isGroup() {
      return this.event.type === eventTypes.GROUP
    },
    isSignUp() {
      return this.event.isSignUpForm
    },
    respondents() {
      return Object.values(this.parsedResponses)
        .map((r) => r.user)
        .filter(Boolean)
    },
    selectedGuestRespondent() {
      if (this.guestAddedAvailability) return this.guestName

      if (this.curRespondents.length !== 1) return ""

      const user = this.parsedResponses[this.curRespondents[0]].user
      return this.isGuest(user) ? user._id : ""
    },
    scheduledEventStyle() {
      const style = {}
      let top, height
      if (this.dragging) {
        top = this.dragStart.row
        height = this.dragCur.row - this.dragStart.row + 1
      } else {
        top = this.curScheduledEvent.hoursOffset * 4
        height = this.curScheduledEvent.hoursLength * 4
      }
      style.top = `calc(${top} * 1rem)`
      style.height = `calc(${height} * 1rem)`
      return style
    },
    signUpBlockBeingDraggedStyle() {
      const style = {}
      let top = 0,
        height = 0
      if (this.dragging) {
        top = this.dragStart.row
        height = this.dragCur.row - this.dragStart.row + 1
      }
      style.top = `calc(${top} * 1rem)`
      style.height = `calc(${height} * 1rem)`
      return style
    },
    /** Parses the responses to the Schej, makes necessary changes based on the type of event, and returns it */
    parsedResponses() {
      const parsed = {}

      // Return calendar availability if group
      if (this.event.type === eventTypes.GROUP) {
        for (const userId in this.event.responses) {
          const calendarEventsByDay = this.groupCalendarEventsByDay[userId]
          if (calendarEventsByDay) {
            // Get manual availability and convert to DOW dates
            const fetchedManualAvailability = this.getManualAvailabilityDow(
              this.fetchedResponses[userId]?.manualAvailability
            )
            const curManualAvailability =
              userId === this.authUser._id
                ? this.getManualAvailabilityDow(this.manualAvailability)
                : {}

            // Get availability from calendar events and use manual availability on the
            // "touched" days
            const availability = this.getAvailabilityFromCalendarEvents({
              calendarEventsByDay,
              includeTouchedAvailability: true,
              fetchedManualAvailability: fetchedManualAvailability ?? {},
              curManualAvailability: curManualAvailability ?? {},
              calendarOptions:
                userId === this.authUser._id
                  ? {
                      bufferTime: this.bufferTime,
                      workingHours: this.workingHours,
                    }
                  : this.fetchedResponses[userId]?.calendarOptions ?? undefined,
            })

            parsed[userId] = {
              ...this.event.responses[userId],
              availability: availability,
            }
          } else {
            parsed[userId] = {
              ...this.event.responses[userId],
              availability: new Set(),
            }
          }
        }
        return parsed
      }

      // Return only current user availability if using blind availabilities and user is not owner
      if (this.event.blindAvailabilityEnabled && !this.isOwner) {
        const guestName = localStorage[this.guestNameKey]
        const userId = this.authUser?._id ?? guestName
        if (userId in this.event.responses) {
          const user = {
            ...this.event.responses[userId].user,
            _id: userId,
          }
          parsed[userId] = {
            ...this.event.responses[userId],
            availability: new Set(
              this.fetchedResponses[userId]?.availability?.map((a) =>
                new Date(a).getTime()
              )
            ),
            ifNeeded: new Set(
              this.fetchedResponses[userId]?.ifNeeded?.map((a) =>
                new Date(a).getTime()
              )
            ),
            user: user,
          }
        }
        return parsed
      }

      // Otherwise, parse responses so that if _id is null (i.e. guest user), then it is set to the guest user's name
      for (const k of Object.keys(this.event.responses)) {
        const newUser = {
          ...this.event.responses[k].user,
          _id: k,
        }
        parsed[k] = {
          ...this.event.responses[k],
          availability: new Set(
            this.fetchedResponses[k]?.availability?.map((a) =>
              new Date(a).getTime()
            )
          ),
          ifNeeded: new Set(
            this.fetchedResponses[k]?.ifNeeded?.map((a) =>
              new Date(a).getTime()
            )
          ),
          user: newUser,
        }
      }
      return parsed
    },
    max() {
      let max = 0
      for (const [dateTime, availability] of this.responsesFormatted) {
        if (availability.size > max) {
          max = availability.size
        }
      }

      return max
    },
    times() {
      /* Returns the times that are encompassed by startTime and endTime */
      const times = []

      for (let i = 0; i < this.event.duration; ++i) {
        const utcTimeNum = this.event.startTime + i
        const localTimeNum = utcTimeToLocalTime(utcTimeNum, this.timezoneOffset)

        times.push({
          hoursOffset: i,
          text: timeNumToTimeText(
            localTimeNum,
            this.timeType === timeTypes.HOUR12
          ),
        })
        times.push({
          hoursOffset: i + 0.25,
        })
        times.push({
          hoursOffset: i + 0.5,
        })
        times.push({
          hoursOffset: i + 0.75,
        })
      }

      return times
    },
    timezoneOffset() {
      if (!("offset" in this.curTimezone)) {
        return new Date().getTimezoneOffset()
      }

      if (this.event.type === eventTypes.DOW) {
        return this.curTimezone.offset * -1
      }

      // Can't just get the offset directly from curTimezone because it doesn't account for dates in the future
      // when daylight savings might be in or out of effect, so instead, we get the timezone for the first date
      // of the event
      return (
        dayjs(this.event.dates[0]).tz(this.curTimezone.value).utcOffset() * -1 // Multiply by -1 because offset is flipped
      )
    },
    userHasResponded() {
      return this.authUser && this.authUser._id in this.parsedResponses
    },
    showLeftZigZag() {
      return this.calendarScrollLeft > 0
    },
    showRightZigZag() {
      return Math.ceil(this.calendarScrollLeft) < this.calendarMaxScroll
    },
    maxDaysPerPage() {
      return this.isPhone ? this.mobileNumDays : 7
    },
    hasNextPage() {
      if (this.event.daysOnly) {
        const lastDay = new Date(this.event.dates[this.event.dates.length - 1])
        const curDate = new Date(this.event.dates[0])
        const monthIndex = curDate.getUTCMonth() + this.page
        const year = curDate.getUTCFullYear()

        const lastDayOfCurMonth = new Date(Date.UTC(year, monthIndex + 1, 0))

        return lastDayOfCurMonth.getTime() < lastDay.getTime()
      }

      return (
        this.event.dates.length - (this.page + 1) * this.maxDaysPerPage > 0 ||
        this.event.type === eventTypes.GROUP
      )
    },
    hasPrevPage() {
      return this.page > 0 || this.event.type === eventTypes.GROUP
    },
    /** Returns whether the event has more than one page */
    hasPages() {
      return this.hasNextPage || this.hasPrevPage
    },

    showStickyRespondents() {
      return (
        this.isPhone &&
        !this.scrolledToRespondents &&
        (this.curTimeslot.row !== -1 ||
          this.curRespondent.length > 0 ||
          this.curRespondents.length > 0)
      )
    },

    // Hint stuff
    hintText() {
      if (this.isPhone) {
        switch (this.state) {
          case this.isGroup && this.states.EDIT_AVAILABILITY:
            return "Toggle which calendars are used. Tap and drag to edit your availability."
          case this.states.EDIT_AVAILABILITY:
            const daysOrTimes = this.event.daysOnly ? "days" : "times"
            if (this.availabilityType === availabilityTypes.IF_NEEDED) {
              return `Tap and drag to add your "if needed" ${daysOrTimes} in yellow.`
            }
            return `Tap and drag to add your "available" ${daysOrTimes} in green.`
          case this.states.SCHEDULE_EVENT:
            return "Tap and drag on the calendar to schedule a Google Calendar event during those times."
          default:
            return ""
        }
      }

      switch (this.state) {
        case this.isGroup && this.states.EDIT_AVAILABILITY:
          return "Toggle which calendars are used. Click and drag to edit your availability."
        case this.states.EDIT_AVAILABILITY:
          const daysOrTimes = this.event.daysOnly ? "days" : "times"
          if (this.availabilityType === availabilityTypes.IF_NEEDED) {
            return `Click and drag to add your "if needed" ${daysOrTimes} in yellow.`
          }
          return `Click and drag to add your "available" ${daysOrTimes} in green.`
        case this.states.SCHEDULE_EVENT:
          return "Click and drag on the calendar to schedule a Google Calendar event during those times."
        default:
          return ""
      }
    },
    hintClosed() {
      return !this.hintState || localStorage[this.hintStateLocalStorageKey]
    },
    hintStateLocalStorageKey() {
      return `closedHintText${this.state}` + ("&isGroup" ? this.isGroup : "")
    },
    hintTextShown() {
      return this.showHintText && this.hintText != "" && !this.hintClosed
    },

    timeslotClassStyle() {
      const classStyles = []
      for (let d = 0; d < this.days.length; ++d) {
        const day = this.days[d]
        for (let t = 0; t < this.times.length; ++t) {
          const time = this.times[t]
          classStyles.push(this.getTimeTimeslotClassStyle(day, time, d, t))
        }
      }
      return classStyles
    },
    dayTimeslotClassStyle() {
      const classStyles = []
      for (let i = 0; i < this.monthDays.length; ++i) {
        classStyles.push(
          this.getDayTimeslotClassStyle(this.monthDays[i].dateObject, i)
        )
      }
      return classStyles
    },
    timeslotVon() {
      const vons = []
      for (let d = 0; d < this.days.length; ++d) {
        for (let t = 0; t < this.times.length; ++t) {
          vons.push(this.getTimeslotVon(t, d))
        }
      }
      return vons
    },
    dayTimeslotVon() {
      const vons = []
      for (let i = 0; i < this.monthDays.length; ++i) {
        const row = Math.floor(i / 7)
        const col = i % 7
        vons.push(this.getTimeslotVon(row, col))
      }
      return vons
    },

    /** Whether to show spinner on top of availability grid */
    showLoader() {
      return (
        // Loading calendar events
        ((this.isGroup || this.alwaysShowCalendarEvents || this.editing) &&
          this.loadingCalendarEvents) ||
        // Loading responses
        this.loadingResponses.loading
      )
    },

    /** Localstorage key containing the guest's name */
    guestNameKey() {
      return `${this.event._id}.guestName`
    },
    /** The guest name stored in localstorage */
    guestName() {
      return localStorage[this.guestNameKey]
    },
    /** Whether a guest has added their availability (saved in localstorage) */
    guestAddedAvailability() {
      return (
        this.guestName?.length > 0 && this.guestName in this.parsedResponses
      )
    },

    /** Returns an array of time blocks representing the current user's availability
     * (used for displaying current user's availability on top of everybody else's availability)
     */
    overlaidAvailability() {
      const overlaidAvailability = []
      this.days.forEach((day, d) => {
        overlaidAvailability.push([])
        let curBlockIndex = 0
        this.times.forEach((time, t) => {
          const date = getDateHoursOffset(day.dateObject, time.hoursOffset)
          const dragAdd =
            this.dragging &&
            this.inDragRange(t, d) &&
            this.dragType === this.DRAG_TYPES.ADD
          const dragRemove =
            this.dragging &&
            this.inDragRange(t, d) &&
            this.dragType === this.DRAG_TYPES.REMOVE

          // Check if timeslot is available or if needed or in the drag region
          if (
            dragAdd ||
            (!dragRemove &&
              (this.availability.has(date.getTime()) ||
                this.ifNeeded.has(date.getTime())))
          ) {
            // Determine whether to render as available or if needed block
            let type = availabilityTypes.AVAILABLE
            if (dragAdd) {
              type = this.availabilityType
            } else {
              type = this.availability.has(date.getTime())
                ? availabilityTypes.AVAILABLE
                : availabilityTypes.IF_NEEDED
            }

            if (curBlockIndex in overlaidAvailability[d]) {
              if (overlaidAvailability[d][curBlockIndex].type === type) {
                // Increase block length if matching type and curBlockIndex exists
                overlaidAvailability[d][curBlockIndex].hoursLength += 0.25
              } else {
                // Add a new block because type is different
                overlaidAvailability[d].push({
                  hoursOffset: time.hoursOffset,
                  hoursLength: 0.25,
                  type,
                })
                curBlockIndex++
              }
            } else {
              // Add a new block because block doesn't exist for current index
              overlaidAvailability[d].push({
                hoursOffset: time.hoursOffset,
                hoursLength: 0.25,
                type,
              })
            }
          } else if (curBlockIndex in overlaidAvailability[d]) {
            // Only increment cur block index if block already exists at the current index
            curBlockIndex++
          }
        })
      })
      return overlaidAvailability
    },

    // Options
    showOverlayAvailabilityToggle() {
      return this.respondents.length > 0 && this.overlayAvailabilitiesEnabled
    },
    showCalendarOptions() {
      return (
        !this.addingAvailabilityAsGuest &&
        this.calendarPermissionGranted &&
        (this.isGroup || (!this.isGroup && !this.userHasResponded))
      )
    },
  },
  methods: {
    ...mapMutations(["setAuthUser"]),
    ...mapActions(["showInfo", "showError"]),

    // -----------------------------------
    //#region Date
    // -----------------------------------

    /** Returns a date object from the dayindex and timeindex given */
    getDateFromDayTimeIndex(dayIndex, timeIndex) {
      return getDateHoursOffset(
        this.days[dayIndex].dateObject,
        this.times[timeIndex].hoursOffset
      )
    },

    /** Returns a date object from the dayindex and hoursoffset given */
    getDateFromDayHoursOffset(dayIndex, hoursOffset) {
      return getDateHoursOffset(this.days[dayIndex].dateObject, hoursOffset)
    },
    //#endregion

    // -----------------------------------
    //#region Respondent
    // -----------------------------------
    mouseOverRespondent(e, id) {
      if (this.curRespondents.length === 0) {
        if (this.state === this.defaultState) {
          this.state = this.states.SINGLE_AVAILABILITY
        }

        this.curRespondent = id
      }
    },
    mouseLeaveRespondent(e) {
      if (this.curRespondents.length === 0) {
        if (this.state === this.states.SINGLE_AVAILABILITY) {
          this.state = this.defaultState
        }

        this.curRespondent = ""
      }
    },
    clickRespondent(e, id) {
      this.state = this.states.SUBSET_AVAILABILITY
      this.curRespondent = ""

      if (this.curRespondentsSet.has(id)) {
        // Remove id
        this.curRespondents = this.curRespondents.filter((r) => r != id)

        // Go back to default state if all users deselected
        if (this.curRespondents.length === 0) {
          this.state = this.defaultState
        }
      } else {
        // Add id
        this.curRespondents.push(id)
      }

      e.stopPropagation()
    },
    deselectRespondents(e) {
      // Don't deselect respondents if toggled best times
      // or if this was fired by clicking on a timeslot
      if (
        e?.target?.previousElementSibling?.id === "show-best-times-toggle" ||
        e?.target?.firstChild?.firstChild?.id === "show-best-times-toggle" ||
        e?.target?.classList?.contains("timeslot") //&& this.isPhone)
      )
        return

      if (this.state === this.states.SUBSET_AVAILABILITY) {
        this.state = this.defaultState
      }

      this.curRespondents = []

      // Stop persisting timeslot
      this.timeslotSelected = false
      this.resetCurTimeslot()
    },

    isGuest(user) {
      return user._id == user.firstName
    },
    //#endregion

    // -----------------------------------
    //#region Aggregate user availability
    // -----------------------------------

    /** Fetches responses from server */
    fetchResponses() {
      if (this.calendarOnly) {
        this.fetchedResponses = this.event.responses
        return
      }

      let timeMin, timeMax
      if (this.event.type === eventTypes.GROUP) {
        if (this.event.dates.length > 0) {
          // Fetch the date range for the current week
          timeMin = new Date(this.event.dates[0])
          timeMax = new Date(this.event.dates[this.event.dates.length - 1])
          timeMax.setDate(timeMax.getDate() + 1)

          // Convert dow dates to discrete dates
          timeMin = dateToDowDate(
            this.event.dates,
            timeMin,
            this.weekOffset,
            true
          )
          timeMax = dateToDowDate(
            this.event.dates,
            timeMax,
            this.weekOffset,
            true
          )
        }
      } else {
        if (this.allDays.length > 0) {
          // Fetch the entire time range of availabilities
          timeMin = new Date(this.allDays[0].dateObject)
          timeMax = new Date(this.allDays[this.allDays.length - 1].dateObject)
          timeMax.setDate(timeMax.getDate() + 1)
        }
      }

      if (!timeMin || !timeMax) return

      // Fetch responses between timeMin and timeMax
      const url = `/events/${
        this.event._id
      }/responses?timeMin=${timeMin.toISOString()}&timeMax=${timeMax.toISOString()}`
      get(url)
        .then((responses) => {
          this.fetchedResponses = responses
          this.getResponsesFormatted()
        })
        .catch((err) => {
          this.showError(
            "There was an error fetching availability! Please refresh the page."
          )
        })
    },
    /** Formats the responses in a map where date/time is mapped to the people that are available then */
    getResponsesFormatted() {
      const lastFetched = new Date().getTime()
      this.loadingResponses.loading = true
      this.loadingResponses.lastFetched = lastFetched

      this.$worker
        .run(
          (days, times, parsedResponses, daysOnly, hideIfNeeded) => {
            // Define functions locally because we can't import functions
            const splitTimeNum = (timeNum) => {
              const hours = Math.floor(timeNum)
              const minutes = Math.floor((timeNum - hours) * 60)
              return { hours, minutes }
            }
            const getDateHoursOffset = (date, hoursOffset) => {
              const { hours, minutes } = splitTimeNum(hoursOffset)
              const newDate = new Date(date)
              newDate.setHours(newDate.getHours() + hours)
              newDate.setMinutes(newDate.getMinutes() + minutes)
              return newDate
            }

            // Create array of all dates in the event
            const dates = []
            if (daysOnly) {
              for (const day of days) {
                dates.push(day.dateObject)
              }
            } else {
              for (const day of days) {
                for (const time of times) {
                  // Iterate through all the times
                  const date = getDateHoursOffset(
                    day.dateObject,
                    time.hoursOffset
                  )
                  dates.push(date)
                }
              }
            }

            // Create a map mapping time to the respondents available during that time
            const formatted = new Map()
            for (const date of dates) {
              formatted.set(date.getTime(), new Set())

              // Check every response and see if they are available for the given time
              for (const response of Object.values(parsedResponses)) {
                // Check availability array
                if (
                  response.availability?.has(date.getTime()) ||
                  (response.ifNeeded?.has(date.getTime()) && !hideIfNeeded)
                ) {
                  formatted.get(date.getTime()).add(response.user._id)
                  continue
                }
              }
            }
            return formatted
          },
          [
            this.allDays,
            this.times,
            this.parsedResponses,
            this.event.daysOnly,
            this.hideIfNeeded,
          ]
        )
        .then((formatted) => {
          // Only set responses formatted for the latest request
          if (lastFetched >= this.loadingResponses.lastFetched) {
            this.responsesFormatted = formatted
          }
        })
        .finally(() => {
          if (this.loadingResponses.lastFetched === lastFetched) {
            this.loadingResponses.loading = false
          }
        })
    },
    /** Returns a set of respondents for the given date/time */
    getRespondentsForHoursOffset(date, hoursOffset) {
      const d = getDateHoursOffset(date, hoursOffset)
      return this.responsesFormatted.get(d.getTime()) ?? new Set()
    },
    showAvailability(row, col) {
      if (this.state === this.states.EDIT_AVAILABILITY && this.isPhone) {
        // Don't show currently selected timeslot when on phone and editing
        return
      }

      // Update current timeslot (the timeslot that has a dotted border around it)
      this.curTimeslot = { row, col }

      if (this.state === this.states.EDIT_AVAILABILITY || this.curRespondent) {
        // Don't show availability when editing or when respondent is selected
        return
      }

      // Update current timeslot availability to show who is available for the given timeslot
      const available =
        this.responsesFormatted.get(
          this.getDateFromRowCol(row, col).getTime()
        ) ?? new Set()
      for (const respondent of this.respondents) {
        if (available.has(respondent._id)) {
          this.curTimeslotAvailability[respondent._id] = true
        } else {
          this.curTimeslotAvailability[respondent._id] = false
        }
      }
    },
    //#endregion

    // -----------------------------------
    //#region Current user availability
    // -----------------------------------
    async refreshAuthUser() {
      this.hasRefreshedAuthUser = true
      await get("/user/profile").then((authUser) => {
        this.setAuthUser(authUser)
      })
    },
    /** resets cur user availability to the response stored on the server */
    resetCurUserAvailability() {
      if (this.event.type === eventTypes.GROUP) {
        this.initSharedCalendarAccounts()
        this.manualAvailability = {}
      }

      this.availability = new Set()
      this.ifNeeded = new Set()
      if (this.userHasResponded) {
        this.populateUserAvailability(this.authUser._id)
      }
    },
    /** Populates the availability set for the auth user from the responses object stored on the server */
    populateUserAvailability(id) {
      this.availability =
        new Set(this.parsedResponses[id]?.availability) ?? new Set()
      this.ifNeeded = new Set(this.parsedResponses[id]?.ifNeeded) ?? new Set()
      this.$nextTick(() => (this.unsavedChanges = false))
    },
    /** Returns a set containing the available times based on the given calendar events object */
    getAvailabilityFromCalendarEvents({
      calendarEventsByDay = [],
      includeTouchedAvailability = false, // Whether to include manual availability for touched days
      fetchedManualAvailability = {}, // Object mapping unix timestamp to array of manual availability (fetched from server)
      curManualAvailability = {}, // Manual availability with edits (takes precedence over fetchedManualAvailability)
      calendarOptions = calendarOptionsDefaults, // User id of the user we are getting availability for
    }) {
      const availability = new Set()
      for (let i = 0; i < this.event.dates.length; ++i) {
        const date = new Date(this.event.dates[i])

        if (includeTouchedAvailability) {
          const endDate = getDateHoursOffset(date, this.event.duration)

          // Check if manual availability has been added for the current date
          let manualAvailabilityAdded = false

          for (const time in curManualAvailability) {
            if (date.getTime() <= time && time <= endDate.getTime()) {
              curManualAvailability[time].forEach((a) => {
                availability.add(new Date(a).getTime())
              })
              delete curManualAvailability[time]
              manualAvailabilityAdded = true
              break
            }
          }

          if (manualAvailabilityAdded) continue

          for (const time in fetchedManualAvailability) {
            if (date.getTime() <= time && time <= endDate.getTime()) {
              fetchedManualAvailability[time].forEach((a) => {
                availability.add(new Date(a).getTime())
              })
              delete fetchedManualAvailability[time]
              manualAvailabilityAdded = true
              break
            }
          }

          if (manualAvailabilityAdded) continue
        }

        // Calculate buffer time
        const bufferTimeInMS = calendarOptions.bufferTime.enabled
          ? calendarOptions.bufferTime.time * 1000 * 60
          : 0

        // Calculate working hours
        const startTimeString = timeNumToTimeString(
          calendarOptions.workingHours.startTime
        )
        const day = getISODateString(getDateWithTimezone(date), true)
        const workingHoursStartDate = dayjs
          .tz(`${day} ${startTimeString}`, this.curTimezone.value)
          .toDate()
        let duration =
          calendarOptions.workingHours.endTime -
          calendarOptions.workingHours.startTime
        if (duration <= 0) duration += 24
        const workingHoursEndDate = getDateHoursOffset(
          workingHoursStartDate,
          duration
        )

        for (const time of this.times) {
          // Check if there exists a calendar event that overlaps [time, time+0.25]
          let startDate = getDateHoursOffset(date, time.hoursOffset)
          const endDate = getDateHoursOffset(date, time.hoursOffset + 0.25)

          // Working hours
          if (calendarOptions.workingHours.enabled) {
            if (
              endDate.getTime() <= workingHoursStartDate.getTime() ||
              startDate.getTime() >= workingHoursEndDate.getTime()
            ) {
              continue
            }
          }

          const index = calendarEventsByDay[i].findIndex((e) => {
            const startDateBuffered = new Date(
              e.startDate.getTime() - bufferTimeInMS
            )
            const endDateBuffered = new Date(
              e.endDate.getTime() + bufferTimeInMS
            )

            const notIntersect =
              dateCompare(endDate, startDateBuffered) <= 0 ||
              dateCompare(startDate, endDateBuffered) >= 0
            return !notIntersect && !e.free
          })
          if (index === -1) {
            availability.add(startDate.getTime())
          }
        }
      }
      return availability
    },
    /** Constructs the availability array using calendarEvents array */
    setAvailabilityAutomatically() {
      // This is not a computed property because we should be able to change it manually from what it automatically fills in
      this.availability = new Set()
      const tmpAvailability = this.getAvailabilityFromCalendarEvents({
        calendarEventsByDay: this.calendarEventsByDay,
        calendarOptions: {
          bufferTime: this.bufferTime,
          workingHours: this.workingHours,
        },
      })

      const pageStartDate = getDateDayOffset(
        new Date(this.event.dates[0]),
        this.page * this.maxDaysPerPage
      )
      const pageEndDate = getDateDayOffset(pageStartDate, this.maxDaysPerPage)
      this.animateAvailability(tmpAvailability, pageStartDate, pageEndDate)
    },
    /** Animate the filling out of availability using setTimeout, between startDate and endDate */
    animateAvailability(availability, startDate, endDate) {
      this.availabilityAnimEnabled = true
      this.availabilityAnimTimeouts = []

      let msPerGroup = 25
      let blocksPerGroup = 2
      if (
        (availability.size / blocksPerGroup) * msPerGroup >
        this.maxAnimTime
      ) {
        blocksPerGroup = (availability.size * msPerGroup) / this.maxAnimTime
      }
      let availabilityArray = [...availability]
      availabilityArray = availabilityArray.filter((a) =>
        isDateBetween(a, startDate, endDate)
      )

      for (let i = 0; i < availabilityArray.length / blocksPerGroup + 1; ++i) {
        const timeout = setTimeout(() => {
          for (const a of availabilityArray.slice(
            i * blocksPerGroup,
            i * blocksPerGroup + blocksPerGroup
          )) {
            this.availability.add(a)
          }
          this.availability = new Set(this.availability)
          if (i >= availabilityArray.length / blocksPerGroup) {
            // Make sure the entire availability has been added (will not be guaranteed when only animating a portion of availability)
            this.availability = new Set(availability)
            this.availabilityAnimTimeouts.push(
              setTimeout(() => {
                this.availabilityAnimEnabled = false

                if (this.showSnackbar) {
                  this.showInfo("Your availability has been autofilled!")
                }
                this.unsavedChanges = false
              }, 500)
            )
          }
        }, i * msPerGroup)

        this.availabilityAnimTimeouts.push(timeout)
      }
    },
    stopAvailabilityAnim() {
      for (const timeout of this.availabilityAnimTimeouts) {
        clearTimeout(timeout)
      }
      this.availabilityAnimEnabled = false
    },
    async submitAvailability(guestPayload = { name: "", email: "" }) {
      let payload = {}

      let type = ""
      // If this is a group submit enabled calendars, otherwise submit availability
      if (this.isGroup) {
        type = "group availability and calendars"
        payload = generateEnabledCalendarsPayload(this.sharedCalendarAccounts)
        payload.manualAvailability = {}
        for (const day of Object.keys(this.manualAvailability)) {
          payload.manualAvailability[day] = [
            ...this.manualAvailability[day],
          ].map((a) => new Date(a))
        }
        payload.calendarOptions = {
          bufferTime: this.bufferTime,
          workingHours: this.workingHours,
        }
      } else {
        type = "availability"
        payload.availability = this.availabilityArray
        payload.ifNeeded = this.ifNeededArray
        if (this.authUser && !this.addingAvailabilityAsGuest) {
          payload.guest = false
        } else {
          payload.guest = true
          payload.name = guestPayload.name
          payload.email = guestPayload.email
          localStorage[this.guestNameKey] = guestPayload.name
        }
      }

      await post(`/events/${this.event._id}/response`, payload)

      // Update analytics
      const addedIfNeededTimes = this.ifNeededArray.length > 0
      if (this.authUser) {
        if (this.authUser._id in this.parsedResponses) {
          this.$posthog?.capture(`Edited ${type}`, {
            eventId: this.event._id,
            addedIfNeededTimes,
          })
        } else {
          this.$posthog?.capture(`Added ${type}`, {
            eventId: this.event._id,
            addedIfNeededTimes,
            // bufferTime: this.bufferTime,
            bufferTime: this.bufferTime.time,
            bufferTimeActive: this.bufferTime.enabled,
            workingHoursEnabled: this.workingHours.enabled,
            workingHoursStartTime: this.workingHours.startTime,
            workingHoursEndTime: this.workingHours.endTime,
          })
        }
      } else {
        if (guestPayload.name in this.parsedResponses) {
          this.$posthog?.capture(`Edited ${type} as guest`, {
            eventId: this.event._id,
            addedIfNeededTimes,
          })
        } else {
          this.$posthog?.capture(`Added ${type} as guest`, {
            eventId: this.event._id,
            addedIfNeededTimes,
          })
        }
      }

      this.refreshEvent()
      this.unsavedChanges = false
    },
    async submitNewSignUpBlocks() {
      if (
        this.signUpBlocksToAddByDay.flat().length +
          this.signUpBlocksByDay.flat().length ===
        0
      ) {
        this.showError("Please add at least one sign-up block!")
        return false
      }

      for (let i = 0; i < this.signUpBlocksToAddByDay.length; ++i) {
        this.signUpBlocksByDay[i] = this.signUpBlocksByDay[i].concat(
          this.signUpBlocksToAddByDay[i]
        )
        this.signUpBlocksToAddByDay[i] = []
      }

      const payload = {
        name: this.event.name,
        duration: this.event.duration,
        dates: this.event.dates,
        type: this.event.type,
        signUpBlocks: this.signUpBlocksByDay.flat().map((block) => {
          return {
            _id: block._id,
            name: block.name,
            capacity: block.capacity,
            startDate: block.startDate,
            endDate: block.endDate,
          }
        }),
      }

      put(`/events/${this.event._id}`, payload)
        .then(() => {
          // window.location.reload()
        })
        .catch((err) => {
          console.log(err)
          this.showError(
            "There was a problem editing this event! Please try again later."
          )
        })

      return true
    },

    async deleteAvailability(name = "") {
      const payload = {}
      if (this.authUser && !this.addingAvailabilityAsGuest) {
        payload.guest = false
        payload.userId = this.authUser._id

        this.$posthog?.capture("Deleted availability", {
          eventId: this.event._id,
        })
      } else {
        payload.guest = true
        payload.name = name

        this.$posthog?.capture("Deleted availability as guest", {
          eventId: this.event._id,
          name,
        })
      }
      await _delete(`/events/${this.event._id}/response`, payload)
      this.availability = new Set()
      if (this.isGroup) this.$router.replace({ name: "home" })
      else this.refreshEvent()
    },
    //#endregion

    // -----------------------------------
    //#region Timeslot
    // -----------------------------------
    setTimeslotSize() {
      /* Gets the dimensions of each timeslot and assigns it to the timeslot variable */
      const timeslotEl = document.querySelector(".timeslot")
      if (timeslotEl) {
        ;({ width: this.timeslot.width, height: this.timeslot.height } =
          timeslotEl.getBoundingClientRect())
      }
    },
    /** Returns a class string and style object for the given time timeslot div */
    getTimeTimeslotClassStyle(day, time, d, t) {
      const date = getDateHoursOffset(day.dateObject, time.hoursOffset)
      const row = t
      const col = d
      const classStyle = this.getTimeslotClassStyle(date, row, col)

      // Add time timeslot specific stuff

      // Animation
      if (this.animateTimeslotAlways || this.availabilityAnimEnabled) {
        classStyle.class += "animate-bg-color "
      }

      // Border style
      if (
        (this.respondents.length > 0 || this.editing) &&
        this.curTimeslot.row === row &&
        this.curTimeslot.col === col
      ) {
        // Dashed border for currently selected timeslot
        classStyle.class +=
          "tw-border tw-border-dashed tw-border-black tw-z-10 "
      } else {
        // Normal border
        const fractionalTime = time.hoursOffset - parseInt(time.hoursOffset)
        if (fractionalTime === 0.25) {
          classStyle.class += "tw-border-b "
          classStyle.style.borderBottomStyle = "dashed"
        } else if (fractionalTime === 0.75) {
          classStyle.class += "tw-border-b "
        }

        classStyle.class += "tw-border-r "
        if (col === 0) classStyle.class += "tw-border-l tw-border-l-gray "
        if (col === this.days.length - 1)
          classStyle.class += "tw-border-r-gray "
        if (row === 0) classStyle.class += "tw-border-t tw-border-t-gray "
        if (row === this.times.length - 1)
          classStyle.class += "tw-border-b tw-border-b-gray "

        if (this.state === this.states.EDIT_AVAILABILITY) {
          classStyle.class += "tw-border-[#999999] "
        } else {
          classStyle.class += "tw-border-[#DDDDDD99] "
        }
      }

      // Change default red:
      if (classStyle.style.backgroundColor === "#E523230D") {
        classStyle.style.backgroundColor = "#E5232333"
      }

      return classStyle
    },
    /** Returns the shared class string and style object for the given timeslot (either time timeslot or day timeslot) */
    getTimeslotClassStyle(date, row, col) {
      let c = ""
      const s = {}

      const timeslotRespondents =
        this.responsesFormatted.get(date.getTime()) ?? new Set()

      // Fill style

      if (this.isSignUp) {
        c += "tw-bg-light-gray "
        return { class: c, style: s }
      }

      if (
        !this.overlayAvailability &&
        this.state === this.states.EDIT_AVAILABILITY
      ) {
        // Set default background color to red (unavailable)
        s.backgroundColor = "#E523230D"

        // Show only current user availability
        const inDragRange = this.inDragRange(row, col)
        if (inDragRange) {
          // Set style if drag range goes over the current timeslot
          if (this.dragType === this.DRAG_TYPES.ADD) {
            if (this.availabilityType === availabilityTypes.AVAILABLE) {
              s.backgroundColor = "#00994C88"
            } else if (this.availabilityType === availabilityTypes.IF_NEEDED) {
              c += "tw-bg-yellow "
            }
          } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
          }
        } else {
          // Otherwise just show the current availability
          // Show current availability from availability set
          if (this.availability.has(date.getTime())) {
            s.backgroundColor = "#00994C88"
          } else if (this.ifNeeded.has(date.getTime())) {
            c += "tw-bg-yellow "
          }
        }
      }

      if (this.state === this.states.SINGLE_AVAILABILITY) {
        // Show only the currently selected respondent's availability
        const respondent = this.curRespondent
        if (timeslotRespondents.has(respondent)) {
          if (this.parsedResponses[respondent]?.ifNeeded?.has(date.getTime())) {
            c += "tw-bg-yellow "
          } else {
            s.backgroundColor = "#00994C88"
          }
        }
      }

      if (
        this.overlayAvailability ||
        this.state === this.states.BEST_TIMES ||
        this.state === this.states.HEATMAP ||
        this.state === this.states.SCHEDULE_EVENT ||
        this.state === this.states.SUBSET_AVAILABILITY
      ) {
        let numRespondents
        let max

        if (
          this.state === this.states.BEST_TIMES ||
          this.state === this.states.HEATMAP ||
          this.state === this.states.SCHEDULE_EVENT
        ) {
          numRespondents = timeslotRespondents.size
          max = this.max
        } else if (this.state === this.states.SUBSET_AVAILABILITY) {
          numRespondents = [...timeslotRespondents].filter((r) =>
            this.curRespondentsSet.has(r)
          ).length

          max = this.curRespondentsMax
        } else if (this.overlayAvailability) {
          if (
            (this.userHasResponded || this.curGuestId?.length > 0) &&
            timeslotRespondents.has(this.authUser?._id ?? this.curGuestId)
          ) {
            // Subtract 1 because we do not want to include current user's availability
            numRespondents = timeslotRespondents.size - 1
            max = this.max
          } else {
            numRespondents = timeslotRespondents.size
            max = this.max
          }
        }

        const totalRespondents =
          this.state === this.states.SUBSET_AVAILABILITY
            ? this.curRespondents.length
            : this.respondents.length

        if (this.defaultState === this.states.BEST_TIMES) {
          if (max > 0 && numRespondents === max) {
            // Only set timeslot to green for the times that most people are available
            if (totalRespondents === 1 || this.overlayAvailability) {
              // Make single responses less saturated
              const green = "#00994CAA"
              s.backgroundColor = green
            } else {
              const green = "#00994C"
              s.backgroundColor = green
            }
          }
        } else if (this.defaultState === this.states.HEATMAP) {
          if (numRespondents > 0) {
            if (totalRespondents === 1) {
              const respondentId =
                this.state === this.states.SUBSET_AVAILABILITY
                  ? this.curRespondents[0]
                  : this.respondents[0]._id
              if (
                this.parsedResponses[respondentId]?.ifNeeded?.has(
                  date.getTime()
                )
              ) {
                c += "tw-bg-yellow "
              } else {
                const green = "#00994CAA"
                s.backgroundColor = green
              }
            } else {
              // Determine color of timeslot based on number of people available
              const frac = numRespondents / max
              const green = "#00994C"
              let alpha
              if (!this.overlayAvailability) {
                alpha = Math.floor(frac * (255 - 30))
                  .toString(16)
                  .toUpperCase()
                  .substring(0, 2)
                  .padStart(2, "0")
                if (
                  frac == 1 &&
                  ((this.curRespondents.length > 0 &&
                    max === this.curRespondents.length) ||
                    (this.curRespondents.length === 0 &&
                      max === this.respondents.length))
                ) {
                  alpha = "FF"
                }
              } else {
                alpha = Math.floor(frac * (255 - 85))
                  .toString(16)
                  .toUpperCase()
                  .substring(0, 2)
                  .padStart(2, "0")
              }

              s.backgroundColor = green + alpha
            }
          }
        }
      }

      return { class: c, style: s }
    },
    getDayTimeslotClassStyle(date, i) {
      const row = Math.floor(i / 7)
      const col = i % 7

      let classStyle
      // Only compute class style for days that are included
      if (this.monthDayIncluded.get(date.getTime())) {
        classStyle = this.getTimeslotClassStyle(date, row, col)
        if (this.state === this.states.EDIT_AVAILABILITY) {
          classStyle.class += "tw-cursor-pointer "
        }

        const backgroundColor = classStyle.style.backgroundColor
        if (
          backgroundColor &&
          lightOrDark(removeTransparencyFromHex(backgroundColor)) === "dark"
        ) {
          classStyle.class += "tw-text-white "
        }
      } else {
        classStyle = {
          class: "tw-bg-off-white tw-text-gray ",
          style: {},
        }
      }

      // Change default red:
      if (classStyle.style.backgroundColor === "#E523230D") {
        classStyle.style.backgroundColor = "#E523233B"
      }

      // Change edit green
      // if (classStyle.style.backgroundColor === "#00994C88") {
      //   classStyle.style.backgroundColor = "#29BC6880"
      // }

      // Border style
      if (
        (this.respondents.length > 0 ||
          this.state === this.states.EDIT_AVAILABILITY) &&
        this.curTimeslot.row === row &&
        this.curTimeslot.col === col &&
        this.monthDayIncluded.get(date.getTime())
      ) {
        // Dashed border for currently selected timeslot
        classStyle.class +=
          "tw-outline-2 tw-outline-dashed tw-outline-black tw-z-10 "
      } else {
        // Normal border
        if (col === 0) classStyle.class += "tw-border-l tw-border-l-gray "
        classStyle.class += "tw-border-r tw-border-r-gray "
        if (col !== 7 - 1) {
          classStyle.style.borderRightStyle = "dashed"
        }

        if (row === 0) classStyle.class += "tw-border-t tw-border-t-gray "
        classStyle.class += "tw-border-b tw-border-b-gray "
        if (row !== Math.floor(this.monthDays.length / 7) - 1) {
          classStyle.style.borderBottomStyle = "dashed"
        }
      }

      return classStyle
    },
    getTimeslotVon(row, col) {
      if (this.interactable) {
        return {
          click: () => {
            if (this.timeslotSelected) {
              // Get rid of persistent timeslot selection if clicked on the same timeslot that is currently being persisted
              if (
                row === this.curTimeslot.row &&
                col === this.curTimeslot.col
              ) {
                this.timeslotSelected = false
              }
            } else if (
              this.state !== this.states.EDIT_AVAILABILITY &&
              (this.userHasResponded || this.guestAddedAvailability)
            ) {
              // Persist timeslot selection if user has already responded
              this.timeslotSelected = true
            }

            this.showAvailability(row, col)
          },
          mousedown: () => {
            // Highlight availability button
            if (
              this.state === this.defaultState &&
              ((!this.isPhone &&
                !(this.userHasResponded || this.guestAddedAvailability)) ||
                this.respondents.length == 0)
            )
              this.highlightAvailabilityBtn()
          },
          mouseover: () => {
            // Only show availability on hover if timeslot is not being persisted
            if (!this.timeslotSelected) {
              this.showAvailability(row, col)
            }
          },
        }
      }
      return {}
    },
    resetCurTimeslot() {
      // Only reset cur timeslot if it isn't being persisted
      if (this.timeslotSelected) return

      this.curTimeslotAvailability = {}
      for (const respondent of this.respondents) {
        this.curTimeslotAvailability[respondent._id] = true
      }
      this.curTimeslot = { row: -1, col: -1 }

      // End drag if mouse left time grid
      this.endDrag()
    },
    //#endregion

    // -----------------------------------
    //#region Editing
    // -----------------------------------
    startEditing() {
      this.state = this.isSignUp
        ? this.states.EDIT_SIGN_UP_BLOCKS
        : this.states.EDIT_AVAILABILITY
      this.availabilityType = availabilityTypes.AVAILABLE
      this.availability = new Set()
      this.ifNeeded = new Set()

      if (this.authUser && !this.addingAvailabilityAsGuest) {
        this.resetCurUserAvailability()
      }
      this.$nextTick(() => (this.unsavedChanges = false))
      this.pageHasChanged = false
    },
    stopEditing() {
      this.state = this.defaultState
      this.stopAvailabilityAnim()

      // Reset options
      this.availabilityType = availabilityTypes.AVAILABLE
      this.overlayAvailability = false
    },
    highlightAvailabilityBtn() {
      this.$emit("highlightAvailabilityBtn")
    },
    editGuestAvailability(id) {
      if (this.authUser) {
        this.$emit("addAvailabilityAsGuest")
      } else {
        this.startEditing()
      }

      this.$nextTick(() => {
        this.populateUserAvailability(id)
        this.$emit("setCurGuestId", id)
      })
    },
    refreshEvent() {
      this.$emit("refreshEvent")
    },
    //#endregion

    // -----------------------------------
    //#region Schedule event
    // -----------------------------------
    scheduleEvent() {
      this.state = this.states.SCHEDULE_EVENT
    },
    cancelScheduleEvent() {
      this.state = this.defaultState
    },

    /** Redirect user to Google Calendar to finish the creation of the event */
    confirmScheduleEvent() {
      if (!this.curScheduledEvent) return
      // Get start date, and end date from the area that the user has dragged out
      const { dayIndex, hoursOffset, hoursLength } = this.curScheduledEvent
      let startDate = this.getDateFromDayHoursOffset(dayIndex, hoursOffset)
      let endDate = this.getDateFromDayHoursOffset(
        dayIndex,
        hoursOffset + hoursLength
      )

      if (this.isWeekly) {
        // Determine offset based on current day of the week.
        // People expect the event to be scheduled in the future, not the past, which is why this check exists
        let offset = 0
        if (new Date().getDay() > startDate.getDay()) {
          offset = 1
        }

        // Transform startDate and endDate to be the current week offset
        startDate = dateToDowDate(this.event.dates, startDate, offset, true)
        endDate = dateToDowDate(this.event.dates, endDate, offset, true)
      }

      // Format email string separated by commas
      const emails = this.respondents.map((r) => {
        // Return email if they are not a guest, otherwise return their name
        if (r.email.length > 0) {
          return r.email
        } else {
          // return `${r.firstName} (no email)`
          return null
        }
      })
      const emailsString = encodeURIComponent(emails.filter(Boolean).join(","))

      // Format start and end date to be in the format required by gcal (remove -, :, and .000)
      const start = startDate.toISOString().replace(/([-:]|\.000)/g, "")
      const end = endDate.toISOString().replace(/([-:]|\.000)/g, "")

      // Construct Google Calendar event creation template url
      const url = `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${encodeURIComponent(
        this.event.name
      )}&dates=${start}/${end}&details=${encodeURIComponent(
        "\n\nThis event was scheduled with schej: https://schej.it/e/"
      )}${this.event._id}&ctz=${this.curTimezone.value}&add=${emailsString}`

      // Navigate to url and reset state
      window.open(url, "_blank")
      this.state = this.defaultState
    },
    //#endregion

    // -----------------------------------
    //#region Drag Stuff
    // -----------------------------------
    normalizeXY(e) {
      /* Normalize the touch event to be relative to element */
      let pageX, pageY
      if ("touches" in e) {
        // is a touch event
        ;({ pageX, pageY } = e.touches[0])
      } else {
        // is a mouse event
        ;({ pageX, pageY } = e)
      }
      const { left, top } = e.currentTarget.getBoundingClientRect()
      const x = pageX - left
      const y = pageY - top - window.scrollY
      return { x, y }
    },
    clampRow(row) {
      if (this.event.daysOnly) {
        row = clamp(row, 0, Math.floor(this.monthDays.length / 7) - 1)
      } else {
        row = clamp(row, 0, this.times.length - 1)
      }
      return row
    },
    clampCol(col) {
      if (this.event.daysOnly) {
        col = clamp(col, 0, 7 - 1)
      } else {
        col = clamp(col, 0, this.days.length - 1)
      }
      return col
    },
    /** Returns row, col for the timeslot we are currently hovering over given the x and y position */
    getRowColFromXY(x, y) {
      const { width, height } = this.timeslot
      let col = Math.floor(x / width)
      let row = Math.floor(y / height)
      row = this.clampRow(row)
      col = this.clampCol(col)
      return {
        row,
        col,
      }
    },
    getDateFromRowCol(row, col) {
      if (this.event.daysOnly) {
        return this.monthDays[row * 7 + col]?.dateObject
      } else {
        if (!this.days[col] || !this.times[row]) return null
        return getDateHoursOffset(
          this.days[col].dateObject,
          this.times[row].hoursOffset
        )
      }
    },
    endDrag() {
      if (!this.allowDrag) return

      if (!this.dragStart || !this.dragCur) return

      // Update availability set based on drag region
      if (this.state === this.states.EDIT_AVAILABILITY) {
        // Determine colInc and rowInc
        let colInc =
          (this.dragCur.col - this.dragStart.col) /
          Math.abs(this.dragCur.col - this.dragStart.col)
        let rowInc =
          (this.dragCur.row - this.dragStart.row) /
          Math.abs(this.dragCur.row - this.dragStart.row)
        if (isNaN(colInc)) colInc = 1
        if (isNaN(rowInc)) rowInc = 1

        // Determine iteration variables
        let rowStart = this.dragStart.row
        let rowMax = this.dragCur.row + rowInc
        let colStart = this.dragStart.col
        let colMax = this.dragCur.col + colInc

        // Correct iteration variables if days only
        if (this.event.daysOnly) {
          colStart = 0
          colMax = 7
          colInc = 1
        }

        // Iterate all selected time slots and either add or remove them
        for (let r = rowStart; r != rowMax; r += rowInc) {
          for (let c = colStart; c != colMax; c += colInc) {
            const date = this.getDateFromRowCol(r, c)

            // Don't add to availability set if month day is not included
            if (
              this.event.daysOnly &&
              (!this.monthDayIncluded.get(date.getTime()) ||
                !this.inDragRange(r, c))
            ) {
              continue
            }

            if (this.dragType === this.DRAG_TYPES.ADD) {
              // Add / remove time from availability set
              if (this.availabilityType === availabilityTypes.AVAILABLE) {
                this.availability.add(date.getTime())
                this.ifNeeded.delete(date.getTime())
              } else if (
                this.availabilityType === availabilityTypes.IF_NEEDED
              ) {
                this.ifNeeded.add(date.getTime())
                this.availability.delete(date.getTime())
              }
            } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
              this.availability.delete(date.getTime())
              this.ifNeeded.delete(date.getTime())
            }

            // Edit manualAvailability set if event is a GROUP
            if (this.event.type === eventTypes.GROUP) {
              const discreteDate = dateToDowDate(
                this.event.dates,
                date,
                this.weekOffset,
                true
              )
              const startDateOfDay = dateToDowDate(
                this.event.dates,
                this.days[c].dateObject,
                this.weekOffset,
                true
              )

              // If date not touched, then add all of the existing calendar availabilities and mark it as touched
              if (!(startDateOfDay.getTime() in this.manualAvailability)) {
                // Create new set
                this.manualAvailability[startDateOfDay.getTime()] = new Set()

                // Add the existing calendar availabilities
                const existingAvailability = this.getAvailabilityForDate(
                  this.days[c].dateObject
                )
                for (const a of existingAvailability) {
                  const convertedDate = dateToDowDate(
                    this.event.dates,
                    new Date(a),
                    this.weekOffset,
                    true
                  )
                  this.manualAvailability[startDateOfDay.getTime()].add(
                    convertedDate.getTime()
                  )
                }
              }

              // Add / remove time from manual availability set
              if (this.dragType === this.DRAG_TYPES.ADD) {
                this.manualAvailability[startDateOfDay.getTime()].add(
                  discreteDate.getTime()
                )
              } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
                this.manualAvailability[startDateOfDay.getTime()].delete(
                  discreteDate.getTime()
                )
              }
            }
          }
        }
        this.availability = new Set(this.availability)
      } else if (this.state === this.states.SCHEDULE_EVENT) {
        // Update scheduled event
        const dayIndex = this.dragStart.col
        const hoursOffset = this.dragStart.row / 4
        const hoursLength = (this.dragCur.row - this.dragStart.row + 1) / 4

        if (hoursLength > 0) {
          this.curScheduledEvent = { dayIndex, hoursOffset, hoursLength }
        } else {
          this.curScheduledEvent = null
        }
      } else if (this.state === this.states.EDIT_SIGN_UP_BLOCKS) {
        // Update sign up blocks
        const dayIndex = this.dragStart.col
        const hoursOffset = this.dragStart.row / 4
        const hoursLength = (this.dragCur.row - this.dragStart.row + 1) / 4
        if (hoursLength > 0) {
          this.signUpBlocksToAddByDay[dayIndex].push(
            this.createSignUpBlock(dayIndex, hoursOffset, hoursLength)
          )
        }
      }

      // Set dragging defaults
      this.dragging = false
      this.dragStart = null
      this.dragCur = null
    },
    inDragRange(row, col) {
      /* Returns whether the given row and col is within the drag range */
      if (this.dragging) {
        if (this.event.daysOnly) {
          if (
            isBetween(row, this.dragStart.row, this.dragCur.row) ||
            isBetween(row, this.dragCur.row, this.dragStart.row)
          ) {
            if (this.dragCur.row < this.dragStart.row) {
              return (
                (this.dragCur.row === row && this.dragCur.col <= col) ||
                (this.dragStart.row === row && this.dragStart.col >= col) ||
                (this.dragStart.row !== row && this.dragCur.row !== row)
              )
            } else if (this.dragCur.row > this.dragStart.row) {
              return (
                (this.dragCur.row === row && this.dragCur.col >= col) ||
                (this.dragStart.row === row && this.dragStart.col <= col) ||
                (this.dragStart.row !== row && this.dragCur.row !== row)
              )
            } else {
              // cur row == start row
              return (
                isBetween(col, this.dragStart.col, this.dragCur.col) ||
                isBetween(col, this.dragCur.col, this.dragStart.col)
              )
            }
          }
          return false
        }

        return (
          (isBetween(row, this.dragStart.row, this.dragCur.row) ||
            isBetween(row, this.dragCur.row, this.dragStart.row)) &&
          (isBetween(col, this.dragStart.col, this.dragCur.col) ||
            isBetween(col, this.dragCur.col, this.dragStart.col))
        )
      }
      return false
    },
    moveDrag(e) {
      if (!this.allowDrag) return
      if (e.touches?.length > 1) return // If dragging with more than one finger
      if (!this.dragStart) return

      e.preventDefault()
      const { row, col } = this.getRowColFromXY(
        ...Object.values(this.normalizeXY(e))
      )

      if (
        this.maxSignUpBlockRowSize &&
        row >= this.dragStart.row + this.maxSignUpBlockRowSize
      ) {
        this.dragCur = {
          row: this.dragStart.row + this.maxSignUpBlockRowSize - 1,
          col,
        }
      } else {
        this.dragCur = { row, col }
      }
    },
    startDrag(e) {
      const { row, col } = this.getRowColFromXY(
        ...Object.values(this.normalizeXY(e))
      )

      // If sign up form, check if trying to drag in a block
      if (this.isSignUp) {
        for (const block of this.signUpBlocksByDay[col].concat(
          this.signUpBlocksToAddByDay[col]
        )) {
          if (
            isBetween(
              row,
              block.hoursOffset * 4,
              (block.hoursOffset + block.hoursLength) * 4 - 1
            )
          ) {
            this.$refs.signUpBlocksList.scrollToSignUpBlock(block._id)
            return
          }
        }
      }

      if (!this.allowDrag) return
      if (e.touches?.length > 1) return // If dragging with more than one finger

      const date = this.getDateFromRowCol(row, col)

      // Dont start dragging if day not included in daysonly event
      if (this.event.daysOnly && !this.monthDayIncluded.get(date.getTime())) {
        return
      }

      this.dragging = true
      this.dragStart = { row, col }
      this.dragCur = { row, col }

      // Prevent scroll
      e.preventDefault()

      // Set drag type
      if (this.isSignUp) {
        this.dragType = this.DRAG_TYPES.ADD
      } else if (
        (this.availabilityType === availabilityTypes.AVAILABLE &&
          this.availability.has(date.getTime())) ||
        (this.availabilityType === availabilityTypes.IF_NEEDED &&
          this.ifNeeded.has(date.getTime()))
      ) {
        this.dragType = this.DRAG_TYPES.REMOVE
      } else {
        this.dragType = this.DRAG_TYPES.ADD
      }
    },
    //#endregion

    // -----------------------------------
    //#region Options
    // -----------------------------------
    getLocalTimezone() {
      const split = new Date(this.event.dates[0])
        .toLocaleTimeString("en-us", { timeZoneName: "short" })
        .split(" ")
      const localTimezone = split[split.length - 1]

      return localTimezone
    },
    onShowBestTimesChange() {
      localStorage["showBestTimes"] = this.showBestTimes
      if (
        this.state == this.states.BEST_TIMES ||
        this.state == this.states.HEATMAP
      )
        this.state = this.defaultState
    },
    toggleShowEditOptions() {
      this.showEditOptions = !this.showEditOptions
      localStorage["showEditOptions"] = this.showEditOptions
    },
    toggleShowEventOptions() {
      this.showEventOptions = !this.showEventOptions
      localStorage["showEventOptions"] = this.showEventOptions
    },
    updateOverlayAvailability(val) {
      this.overlayAvailability = !!val
    },
    //#endregion

    // -----------------------------------
    //#region Scroll
    // -----------------------------------
    onCalendarScroll(e) {
      this.calendarMaxScroll = e.target.scrollWidth - e.target.offsetWidth
      this.calendarScrollLeft = e.target.scrollLeft
    },
    onScroll(e) {
      this.checkElementsVisible()
    },
    /** Checks whether certain elements are visible and sets variables accoringly */
    checkElementsVisible() {
      const optionsSectionEl = this.$refs.optionsSection
      if (optionsSectionEl) {
        this.optionsVisible = isElementInViewport(optionsSectionEl, {
          bottomOffset: -64,
        })
      }

      const respondentsListEl = this.$refs.respondentsList?.$el
      if (respondentsListEl) {
        this.scrolledToRespondents = isElementInViewport(respondentsListEl, {
          bottomOffset: -64,
        })
      }
    },
    //#endregion

    // -----------------------------------
    //#region Pagination
    // -----------------------------------
    nextPage(e) {
      e.stopImmediatePropagation()
      if (this.event.type === eventTypes.GROUP) {
        // Go to next page if there are still more days left to see
        // Otherwise, update week offset
        if ((this.page + 1) * this.maxDaysPerPage < this.event.dates.length) {
          this.page++
        } else {
          this.page = 0
          this.$emit("update:weekOffset", this.weekOffset + 1)
        }
      } else {
        this.page++
      }
      this.pageHasChanged = true
    },
    prevPage(e) {
      e.stopImmediatePropagation()
      if (this.event.type === eventTypes.GROUP) {
        // Go to prev page if there is a prev page
        // Otherwise, update week offset
        if (this.page > 0) {
          this.page--
        } else {
          this.page =
            Math.ceil(this.event.dates.length / this.maxDaysPerPage) - 1
          this.$emit("update:weekOffset", this.weekOffset - 1)
        }
      } else {
        this.page--
      }
      this.pageHasChanged = true
    },
    //#endregion

    // -----------------------------------
    //#region Resize
    // -----------------------------------
    onResize() {
      this.setTimeslotSize()
    },
    //#endregion

    // -----------------------------------
    //#region hint
    // -----------------------------------
    closeHint() {
      this.hintState = false
      localStorage[this.hintStateLocalStorageKey] = true
    },
    //#endregion

    // -----------------------------------
    //#region Group
    // -----------------------------------

    /** Toggles calendar account - in groups to enable/disable calendars */
    toggleCalendarAccount(payload) {
      this.sharedCalendarAccounts[
        getCalendarAccountKey(payload.email, payload.calendarType)
      ].enabled = payload.enabled
      this.sharedCalendarAccounts = JSON.parse(
        JSON.stringify(this.sharedCalendarAccounts)
      )
    },

    /** Toggles sub calendar account - in groups to enable/disable sub calendars */
    toggleSubCalendarAccount(payload) {
      this.sharedCalendarAccounts[
        getCalendarAccountKey(payload.email, payload.calendarType)
      ].subCalendars[payload.subCalendarId].enabled = payload.enabled
      this.sharedCalendarAccounts = JSON.parse(
        JSON.stringify(this.sharedCalendarAccounts)
      )
    },

    /** Sets the initial sharedCalendarAccounts object */
    initSharedCalendarAccounts() {
      if (!this.authUser) return

      // Init shared calendar accounts to current calendar accounts
      this.sharedCalendarAccounts = JSON.parse(
        JSON.stringify(this.authUser.calendarAccounts)
      )

      // Disable all calendars
      for (const id in this.sharedCalendarAccounts) {
        this.sharedCalendarAccounts[id].enabled = false
        if (this.sharedCalendarAccounts[id].subCalendars) {
          for (const subCalendarId in this.sharedCalendarAccounts[id]
            .subCalendars) {
            this.sharedCalendarAccounts[id].subCalendars[
              subCalendarId
            ].enabled = false
          }
        }
      }

      // Enable calendars based on responses
      if (this.authUser._id in this.event.responses) {
        const enabledCalendars =
          this.event.responses[this.authUser._id].enabledCalendars

        for (const id in enabledCalendars) {
          this.sharedCalendarAccounts[id].enabled = true

          enabledCalendars[id].forEach((subCalendarId) => {
            this.sharedCalendarAccounts[id].subCalendars[
              subCalendarId
            ].enabled = true
          })
        }
      }
    },

    /** Based on the date, determine whether it has been touched */
    isTouched(date, availability = [...this.availability]) {
      const start = new Date(date)
      const end = new Date(date)
      end.setHours(end.getHours() + this.event.duration)

      for (const a of availability) {
        const availableTime = new Date(a).getTime()
        if (
          start.getTime() <= availableTime &&
          availableTime <= end.getTime()
        ) {
          return true
        }
      }

      return false
    },

    /** Returns a subset of availability for the current date */
    getAvailabilityForDate(date, availability = [...this.availability]) {
      const start = new Date(date)
      const end = new Date(date)
      end.setHours(end.getHours() + this.event.duration)

      const subset = new Set()
      for (const a of availability) {
        const availableTime = new Date(a).getTime()
        if (
          start.getTime() <= availableTime &&
          availableTime <= end.getTime()
        ) {
          subset.add(availableTime)
        }
      }

      return subset
    },

    /** Returns a copy of the manual availability, converted to dow dates */
    getManualAvailabilityDow(manualAvailability = this.manualAvailability) {
      if (!manualAvailability) return null

      const manualAvailabilityDow = {}
      for (const time in manualAvailability) {
        const dowTime = dateToDowDate(
          this.event.dates,
          new Date(parseInt(time)),
          this.weekOffset
        ).getTime()
        manualAvailabilityDow[dowTime] = [...manualAvailability[time]].map(
          (a) => dateToDowDate(this.event.dates, new Date(a), this.weekOffset)
        )
      }
      return manualAvailabilityDow
    },
    //#endregion

    // -----------------------------------
    //#region Sign up form
    // -----------------------------------

    /** Creates a sign up block for the current day and hour offset */
    createSignUpBlock(dayIndex, hoursOffset, hoursLength) {
      const timeBlock = getTimeBlock(
        this.days[dayIndex].dateObject,
        hoursOffset,
        hoursLength
      )

      return {
        _id: ObjectID().toString(),
        capacity: 1,
        name: this.newSignUpBlockName,
        ...timeBlock,
        hoursOffset,
        hoursLength,
      }
    },

    /** Updates the sign up block with the same id */
    editSignUpBlock(signUpBlock) {
      this.signUpBlocksByDay.forEach((blocksInDay, dayIndex) => {
        blocksInDay.forEach((block, blockIndex) => {
          if (signUpBlock._id === block._id) {
            this.signUpBlocksByDay[dayIndex][blockIndex] = signUpBlock
            this.signUpBlocksByDay = [...this.signUpBlocksByDay]
            return
          }
        })
      })

      this.signUpBlocksToAddByDay.forEach((blocksInDay, dayIndex) => {
        blocksInDay.forEach((block, blockIndex) => {
          if (signUpBlock._id === block._id) {
            this.signUpBlocksToAddByDay[dayIndex][blockIndex] = signUpBlock
            this.signUpBlocksToAddByDay = [...this.signUpBlocksToAddByDay]
            return
          }
        })
      })
    },

    /** Deletes the sign up block with the id */
    deleteSignUpBlock(signUpBlockId) {
      this.signUpBlocksByDay.forEach((blocksInDay, dayIndex) => {
        blocksInDay.forEach((block, blockIndex) => {
          if (signUpBlockId === block._id) {
            this.signUpBlocksByDay[dayIndex].splice(blockIndex, 1)
            return
          }
        })
      })

      this.signUpBlocksToAddByDay.forEach((blocksInDay, dayIndex) => {
        blocksInDay.forEach((block, blockIndex) => {
          if (signUpBlockId === block._id) {
            this.signUpBlocksToAddByDay[dayIndex].splice(blockIndex, 1)
            return
          }
        })
      })
    },

    /** Reloads all the data for the sign up form */
    resetSignUpForm() {
      /** Split sign up blocks by day */
      this.signUpBlocksByDay = splitTimeBlocksByDay(
        this.event,
        this.event.signUpBlocks ?? []
      )

      this.resetSignUpBlocksToAddByDay()

      /** Populate sign up block responses */
      for (const userId in this.event.signUpResponses) {
        const signUpResponse = this.event.signUpResponses[userId]
        for (const signUpBlockId of signUpResponse.signUpBlockIds) {
          const signUpBlock = this.signUpBlocksByDay
            .flat()
            .find((signUpBlock) => signUpBlock._id === signUpBlockId)

          if (!signUpBlock.responses) signUpBlock.responses = []
          signUpBlock.responses.push(signUpResponse)
        }
      }
    },

    /** Initialize sign up blocks to be added array */
    resetSignUpBlocksToAddByDay() {
      this.signUpBlocksToAddByDay = []
      for (const day of this.signUpBlocksByDay) {
        this.signUpBlocksToAddByDay.push([])
      }
    },

    /** Emits sign up for block to parent element */
    handleSignUpBlockClick(block) {
      if (!this.alreadyRespondedToSignUpForm && !this.isOwner)
        this.$emit("signUpForBlock", block)
    },

    //#endregion

    /** Recalculate availability the calendar based on calendar events */
    reanimateAvailability() {
      if (
        this.state === this.states.EDIT_AVAILABILITY &&
        this.authUser &&
        !(this.authUser?._id in this.event.responses) && // User hasn't responded yet
        !this.loadingCalendarEvents &&
        (!this.unsavedChanges || this.availabilityAnimEnabled)
      ) {
        for (const timeout of this.availabilityAnimTimeouts) {
          clearTimeout(timeout)
        }
        this.setAvailabilityAutomatically()
      }
    },
  },
  watch: {
    availability() {
      if (this.state === this.states.EDIT_AVAILABILITY) {
        this.unsavedChanges = true
      }
    },
    event: {
      immediate: true,
      handler() {
        this.initSharedCalendarAccounts()
        this.fetchResponses()
      },
    },
    state(nextState, prevState) {
      this.$nextTick(() => this.checkElementsVisible())

      // Reset scheduled event when exiting schedule event state
      if (prevState === this.states.SCHEDULE_EVENT) {
        this.curScheduledEvent = null
      } else if (prevState === this.states.EDIT_AVAILABILITY) {
        this.unsavedChanges = false
      }
    },
    respondents: {
      immediate: true,
      handler() {
        this.curTimeslotAvailability = {}
        for (const respondent of this.respondents) {
          this.curTimeslotAvailability[respondent._id] = true
        }
      },
    },
    calendarEventsByDay() {
      this.reanimateAvailability()
    },
    page() {
      this.$nextTick(() => {
        this.setTimeslotSize()
      })
    },
    showStickyRespondents: {
      immediate: true,
      handler(cur) {
        clearTimeout(this.delayedShowStickyRespondentsTimeout)
        this.delayedShowStickyRespondentsTimeout = setTimeout(() => {
          this.delayedShowStickyRespondents = cur
        }, 100)
      },
    },
    maxDaysPerPage() {
      // Set page to 0 if user switches from portrait to landscape orientation and we're on an invalid page number,
      // i.e. we're on a page that displays 0 days
      if (this.page * this.maxDaysPerPage >= this.event.dates.length) {
        this.page = 0
      }
    },
    mobileNumDays() {
      // Save mobile num days in localstorage
      localStorage["mobileNumDays"] = this.mobileNumDays

      // Set timeslot size because it has changed
      this.$nextTick(() => {
        this.setTimeslotSize()
      })
    },
    weekOffset() {
      if (this.event.type === eventTypes.GROUP) {
        this.fetchResponses()
      }
    },
    hideIfNeeded() {
      this.getResponsesFormatted()
    },
    parsedResponses() {
      // Theoretically, parsed responses should only be changing for groups
      this.getResponsesFormatted()

      // Repopulate user availability when editing availability (this happens when switching weeks in a group)
      if (
        this.event.type === eventTypes.GROUP &&
        this.state === this.states.EDIT_AVAILABILITY &&
        this.authUser
      ) {
        this.availability = new Set()
        this.populateUserAvailability(this.authUser._id)
      }
    },
    showBestTimes() {
      this.onShowBestTimesChange()
    },
    startCalendarOnMonday() {
      localStorage["startCalendarOnMonday"] = this.startCalendarOnMonday
    },
    bufferTime(cur, prev) {
      if (cur.enabled !== prev.enabled || cur.enabled) {
        this.reanimateAvailability()
      }
    },
    workingHours(cur, prev) {
      if (cur.enabled !== prev.enabled || cur.enabled) {
        this.reanimateAvailability()
      }
    },
    timeType() {
      localStorage["timeType"] = this.timeType
    },
  },
  created() {
    this.resetCurUserAvailability()

    addEventListener("click", this.deselectRespondents)
  },
  mounted() {
    // Set initial state to best_times or heatmap depending on show best times toggle.
    this.state = this.showBestTimes ? "best_times" : "heatmap"

    // Set calendar options defaults
    if (this.authUser) {
      this.bufferTime =
        this.authUser?.calendarOptions?.bufferTime ??
        calendarOptionsDefaults.bufferTime
      this.workingHours =
        this.authUser?.calendarOptions?.workingHours ??
        calendarOptionsDefaults.workingHours
      if (this.isGroup) {
        if (this.event.responses[this.authUser._id]?.calendarOptions) {
          // Update calendar options if user has changed them for this specific group
          const { bufferTime, workingHours } =
            this.event.responses[this.authUser._id]?.calendarOptions
          if (bufferTime) this.bufferTime = bufferTime
          if (workingHours) this.workingHours = workingHours
        } else {
          this.bufferTime = calendarOptionsDefaults.bufferTime
          this.workingHours = calendarOptionsDefaults.workingHours
        }
      }
    }

    // Set initial calendar max scroll
    // this.calendarMaxScroll =
    //   this.$refs.calendar.scrollWidth - this.$refs.calendar.offsetWidth

    // Get timeslot size
    this.setTimeslotSize()
    addEventListener("resize", this.onResize)
    addEventListener("scroll", this.onScroll)
    if (!this.calendarOnly) {
      const timesEl = document.getElementById("drag-section")
      if (isTouchEnabled()) {
        timesEl.addEventListener("touchstart", this.startDrag)
        timesEl.addEventListener("touchmove", this.moveDrag)
        timesEl.addEventListener("touchend", this.endDrag)
        timesEl.addEventListener("touchcancel", this.endDrag)
      }
      timesEl.addEventListener("mousedown", this.startDrag)
      timesEl.addEventListener("mousemove", this.moveDrag)
      timesEl.addEventListener("mouseup", this.endDrag)
    }

    // Parse sign up blocks and responses
    this.resetSignUpForm()
  },
  beforeDestroy() {
    removeEventListener("click", this.deselectRespondents)
    removeEventListener("resize", this.onResize)
    removeEventListener("scroll", this.onScroll)
  },
  components: {
    AlertText,
    AvailabilityTypeToggle,
    ExpandableSection,
    BufferTimeSwitch,
    UserAvatarContent,
    ZigZag,
    ConfirmDetailsDialog,
    ToolRow,
    CalendarAccounts,
    RespondentsList,
    Advertisement,
    GCalWeekSelector,
    WorkingHoursToggle,
    SignUpBlock,
    SignUpCalendarBlock,
    SignUpBlocksList,
  },
}
</script>
