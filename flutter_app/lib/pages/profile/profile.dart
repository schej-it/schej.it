import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/auth_service.dart';
import 'package:flutter_app/pages/profile/expansion_panel.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:provider/provider.dart';
import 'package:url_launcher/url_launcher.dart';

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
  int _visibility = 0;
  bool _noficiations = false;
  double settingHeight = 61.0;
  var options = [
    Setting("Everything", "Other users can view event names of your schej"),
    Setting("Busy/free", "Other users view your events as “busy”"),
    Setting("Invisible", "Other people cannot view your schej"),
  ];

  @override
  void initState() {
    super.initState();

    final api = context.read<ApiService>();
    _visibility = api.authUser?.visibility ?? 0;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(titleString: 'Profile', isRoot: true),
      body: SingleChildScrollView(
        child: Container(
          padding: const EdgeInsets.all(25.0),
          child: Consumer<ApiService>(
            builder: (context, api, child) => Column(
              children: [
                _buildProfileSection(api),
                const SizedBox(height: 40.0),
                _buildSettingsSection(api),
                const SizedBox(height: 150.0),
                _buildSignOutSection(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildProfileSection(ApiService api) {
    return Center(
      child: Column(children: [
        CircleAvatar(
          backgroundImage: NetworkImage('${api.authUser?.picture}'),
          radius: 35,
        ),
        Text('${api.authUser?.firstName} ${api.authUser?.lastName}',
            style: SchejFonts.header),
        Text('${api.authUser?.email}', style: SchejFonts.body),
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

  Widget _buildSettingsSection(ApiService api) {
    return Column(
      children: [
        _buildVisibilitySetting(api),
        const SizedBox(height: 10),
        _buildNotificationsSetting(),
        const SizedBox(height: 10),
        _buildFeedbackSection(),
      ],
    );
  }

  Widget _buildVisibilitySetting(ApiService api) {
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
                              api.updateUserVisibility(_visibility);
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

  Widget _buildFeedbackSection() {
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
          const Text("Feedback", style: SchejFonts.medium),
          Row(
            children: [
              IconButton(
                  onPressed: () => {_launchUrl()},
                  icon: const Icon(MdiIcons.openInNew))
            ],
          ),
        ],
      ),
    );
  }

  Future<void> _launchUrl() async {
    final Uri url = Uri.parse('https://forms.gle/9AgRy4PQfWfVuBnw8');
    if (!await launchUrl(url, mode: LaunchMode.externalApplication)) {
      throw 'Could not launch $url';
    }
  }
}
