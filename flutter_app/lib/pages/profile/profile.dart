import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:flutter_app/pages/profile/expansion_panel.dart';
import 'package:provider/provider.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({Key? key}) : super(key: key);

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class Setting {
  late String name;
  late String description;
  Setting(this.name, this.description);
}

class _ProfilePageState extends State<ProfilePage> {
  bool _active = false;
  int _visibility = 1;
  bool _noficiations = false;
  double settingHeight = 61.0;
  var options = [
    Setting("Everything", "Other users can view event names of your schej"),
    Setting("Busy/free", "Other users view your events as “busy”"),
    Setting("Invisible", "Other people cannot view your schej"),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(titleString: 'Friends', isRoot: true),
      body: Container(
        padding: const EdgeInsets.all(25.0),
        child: SingleChildScrollView(
          child: Column(
            children: [
              _buildProfileSection(),
              const SizedBox(height: 40.0),
              _buildSettingsSection(),
              const SizedBox(height: 150.0),
              _buildSignOutSection(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildProfileSection() {
    return Center(
      child: Column(children: const [
        CircleAvatar(
          backgroundImage: NetworkImage('https://i.imgur.com/HkLY72h.jpg'),
          radius: 35,
        ),
        Text("Tony Xin", style: SchejFonts.header),
        Text("tonyxin@berkeley.edu", style: SchejFonts.body),
      ]),
    );
  }

  Widget _buildSignOutSection() {
    return Center(
      child: TextButton(
        onPressed: () {
          final authService = context.read<AuthService>();
          authService.signOut();
        },
        child: Text('Sign out',
            style: SchejFonts.header.copyWith(color: SchejColors.red)),
      ),
    );
  }

  Widget _buildSettingsSection() {
    return Column(
      children: [
        _buildVisibilitySetting(),
        const SizedBox(height: 10),
        _buildNotificationsSetting(),
      ],
    );
  }

  Widget _buildVisibilitySetting() {
    return CustomExpansionPanelList(
      expansionCallback: (panelIndex, isExpanded) {
        _active = !_active;
        setState(() {});
      },
      key: const Key("visibility"),
      children: <ExpansionPanel>[
        ExpansionPanel(
            headerBuilder: (context, isExpanded) {
              return Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Padding(
                    padding: EdgeInsets.only(left: 20.0),
                    child: Text("Visibility", style: SchejFonts.medium),
                  ),
                  Row(
                    children: [
                      Text(options[_visibility].name, style: SchejFonts.body),
                    ],
                  ),
                ],
              );
            },
            body: Column(
              children: [
                for (int i = 0; i < options.length; i++)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 8.0),
                    child: ListTile(
                      title: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            options[i].name,
                            style: SchejFonts.medium,
                          ),
                          const SizedBox(height: 2.0),
                          Text(
                            options[i].description,
                            style: SchejFonts.body,
                          ),
                        ],
                      ),
                      leading: Transform.scale(
                        scale: 1.4,
                        child: Radio(
                          value: i,
                          groupValue: _visibility,
                          activeColor: SchejColors.green,
                          onChanged: (int? value) {
                            setState(() {
                              _visibility = value!;
                            });
                          },
                        ),
                      ),
                    ),
                  ),
              ],
            ),
            isExpanded: _active,
            canTapOnHeader: true)
      ],
    );
  }

  Widget _buildNotificationsSetting() {
    return Container(
      height: settingHeight,
      padding: const EdgeInsets.only(left: 20.0, right: 20.0),
      decoration: BoxDecoration(
        boxShadow: [SchejConstants.boxShadow],
        color: SchejColors.white,
        borderRadius: const BorderRadius.all(Radius.circular(10.0)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          const Text("Notifications", style: SchejFonts.medium),
          Row(
            children: [
              Switch.adaptive(
                value: _noficiations,
                onChanged: (value) {
                  setState(() {
                    _noficiations = value;
                  });
                },
                activeTrackColor: SchejColors.lightGreen,
                activeColor: SchejColors.green,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
